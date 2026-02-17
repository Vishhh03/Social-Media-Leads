package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/meta"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

// ChannelHandler handles channel management endpoints.
type ChannelHandler struct {
	Store          *store.Storage
	TokenRefresher *meta.TokenRefresher
}

// ConnectChannelRequest is the expected body for connecting a channel.
type ConnectChannelRequest struct {
	Platform    string `json:"platform" binding:"required"`
	AccountID   string `json:"account_id" binding:"required"`
	AccountName string `json:"account_name"`
	AccessToken string `json:"access_token" binding:"required"`
}

// ConnectChannel creates a new channel connection.
func (h *ChannelHandler) ConnectChannel(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ConnectChannelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate platform
	switch req.Platform {
	case "whatsapp", "instagram", "facebook":
		// OK
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Platform must be whatsapp, instagram, or facebook"})
		return
	}

	// Try to exchange for a long-lived token
	accessToken := req.AccessToken
	tokenExpiry := time.Time{}

	if h.TokenRefresher != nil && h.TokenRefresher.AppID != "" {
		tokenResp, err := h.TokenRefresher.ExchangeForLongLivedToken(c.Request.Context(), req.AccessToken)
		if err != nil {
			log.Printf("[Channel] Token exchange failed (using short-lived token): %v", err)
		} else {
			accessToken = tokenResp.AccessToken
			tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
			log.Printf("[Channel] ✅ Exchanged for long-lived token (expires in %d seconds)", tokenResp.ExpiresIn)
		}
	}

	channel := &models.Channel{
		UserID:      userID.(int64),
		Platform:    req.Platform,
		AccountID:   req.AccountID,
		AccountName: req.AccountName,
		AccessToken: accessToken,
		TokenExpiry: tokenExpiry,
		IsActive:    true,
	}

	if err := h.Store.CreateChannel(c.Request.Context(), channel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect channel"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Channel connected",
		"channel": gin.H{
			"id":           channel.ID,
			"platform":     channel.Platform,
			"account_id":   channel.AccountID,
			"account_name": channel.AccountName,
			"is_active":    channel.IsActive,
			"created_at":   channel.CreatedAt,
		},
	})
}

// ListChannels returns all channels for the current user.
func (h *ChannelHandler) ListChannels(c *gin.Context) {
	userID, _ := c.Get("user_id")

	channels, err := h.Store.GetChannelsByUser(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch channels"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"channels": channels,
		"count":    len(channels),
	})
}

// DisconnectChannel soft-deletes a channel.
func (h *ChannelHandler) DisconnectChannel(c *gin.Context) {
	userID, _ := c.Get("user_id")
	channelID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
		return
	}

	if err := h.Store.DeleteChannel(c.Request.Context(), channelID, userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect channel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Channel disconnected"})
}
