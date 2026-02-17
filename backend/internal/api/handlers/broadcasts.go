package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/cache"
	"github.com/social-media-lead/backend/internal/meta"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

// BroadcastHandler handles broadcast messaging endpoints.
type BroadcastHandler struct {
	Store      *store.Storage
	MetaClient *meta.Client
	Redis      *cache.RedisClient
}

// CreateBroadcastRequest is the expected body for creating a broadcast.
type CreateBroadcastRequest struct {
	Name        string     `json:"name" binding:"required"`
	Content     string     `json:"content" binding:"required"`
	MediaURL    string     `json:"media_url"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

// CreateBroadcast creates a new broadcast draft.
func (h *BroadcastHandler) CreateBroadcast(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateBroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	broadcast := &models.Broadcast{
		UserID:      userID.(int64),
		Name:        req.Name,
		Content:     req.Content,
		MediaURL:    req.MediaURL,
		Status:      "draft",
		ScheduledAt: req.ScheduledAt,
	}

	if err := h.Store.CreateBroadcast(c.Request.Context(), broadcast); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create broadcast"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Broadcast created",
		"broadcast": broadcast,
	})
}

// ListBroadcasts returns all broadcasts for the current user.
func (h *BroadcastHandler) ListBroadcasts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	broadcasts, err := h.Store.GetBroadcastsByUser(c.Request.Context(), userID.(int64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch broadcasts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"broadcasts": broadcasts,
		"count":      len(broadcasts),
	})
}

// SendBroadcast executes a broadcast, sending the message to all user contacts.
func (h *BroadcastHandler) SendBroadcast(c *gin.Context) {
	userID, _ := c.Get("user_id")
	broadcastID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid broadcast ID"})
		return
	}

	broadcast, err := h.Store.GetBroadcastByID(c.Request.Context(), broadcastID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Broadcast not found"})
		return
	}

	if broadcast.UserID != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if broadcast.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Broadcast already sent or sending"})
		return
	}

	_ = h.Store.UpdateBroadcastStatus(c.Request.Context(), broadcastID, "sending", 0, 0)

	go h.executeBroadcast(broadcast)

	c.JSON(http.StatusOK, gin.H{
		"message": "Broadcast sending started",
		"status":  "sending",
	})
}

// executeBroadcast sends the broadcast message to all contacts for the user.
func (h *BroadcastHandler) executeBroadcast(broadcast *models.Broadcast) {
	ctx := context.Background()

	contacts, err := h.Store.GetContactsByUser(ctx, broadcast.UserID, 10000, 0)
	if err != nil {
		log.Printf("[Broadcast] Failed to fetch contacts: %v", err)
		_ = h.Store.UpdateBroadcastStatus(ctx, broadcast.ID, "failed", 0, 0)
		return
	}

	channels, err := h.Store.GetChannelsByUser(ctx, broadcast.UserID)
	if err != nil {
		log.Printf("[Broadcast] Failed to fetch channels: %v", err)
		_ = h.Store.UpdateBroadcastStatus(ctx, broadcast.ID, "failed", 0, 0)
		return
	}

	channelMap := make(map[int64]*models.Channel)
	for i := range channels {
		channelMap[channels[i].ID] = &channels[i]
	}

	// Set broadcast dedup expiry (24 hours)
	if h.Redis != nil {
		_ = h.Redis.ExpireBroadcastSet(ctx, broadcast.ID, 24*time.Hour)
	}

	totalSent := 0
	totalFailed := 0

	for _, contact := range contacts {
		// Check Redis dedup to prevent double-sends
		if h.Redis != nil {
			alreadySent, _ := h.Redis.WasBroadcastSent(ctx, broadcast.ID, contact.ID)
			if alreadySent {
				log.Printf("[Broadcast] Skipping contact #%d (already sent)", contact.ID)
				continue
			}
		}

		// Rate limit check
		if h.Redis != nil {
			result, _ := h.Redis.CheckRateLimit(ctx, "meta_api:broadcast", 200, 1*time.Minute)
			if result != nil && !result.Allowed {
				log.Printf("[Broadcast] Rate limited, pausing 60s...")
				time.Sleep(60 * time.Second)
			}
		}

		channel, ok := channelMap[contact.ChannelID]
		if !ok || !channel.IsActive {
			totalFailed++
			continue
		}

		result, err := h.MetaClient.SendMessage(
			contact.Platform,
			channel.AccountID,
			contact.PlatformUserID,
			broadcast.Content,
			channel.AccessToken,
		)

		if err != nil || !result.Success {
			totalFailed++
			log.Printf("[Broadcast] Failed to send to contact #%d: %v", contact.ID, err)
			continue
		}

		// Mark as sent in Redis
		if h.Redis != nil {
			_ = h.Redis.MarkBroadcastSent(ctx, broadcast.ID, contact.ID)
			h.Redis.LogEvent(ctx, "message_sent", broadcast.UserID)
		}

		// Store outbound message
		msg := &models.Message{
			UserID:        broadcast.UserID,
			ChannelID:     channel.ID,
			ContactID:     contact.ID,
			Platform:      contact.Platform,
			Direction:     "outbound",
			Content:       broadcast.Content,
			MessageType:   "text",
			PlatformMsgID: result.MessageID,
			Status:        "sent",
			IsAutomated:   true,
		}
		_ = h.Store.CreateMessage(ctx, msg)

		totalSent++
	}

	_ = h.Store.UpdateBroadcastStatus(ctx, broadcast.ID, "sent", totalSent, totalFailed)
	log.Printf("[Broadcast] ✅ Completed '%s': %d sent, %d failed", broadcast.Name, totalSent, totalFailed)
}
