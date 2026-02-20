package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/meta"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

// InboxHandler handles unified inbox endpoints.
type InboxHandler struct {
	Store      store.Store
	MetaClient *meta.Client
}

// GetConversations returns the last message per contact for the current user (inbox list).
func (h *InboxHandler) GetConversations(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	conversations, err := h.Store.GetConversations(c.Request.Context(), userID.(int64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": conversations,
		"count":         len(conversations),
	})
}

// GetMessages returns messages for a specific contact (chat thread).
func (h *InboxHandler) GetMessages(c *gin.Context) {
	contactID, err := strconv.ParseInt(c.Param("contact_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.Store.GetMessagesByContact(c.Request.Context(), contactID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"count":    len(messages),
	})
}

// GetContacts returns all leads/contacts for the current user.
func (h *InboxHandler) GetContacts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	contacts, err := h.Store.GetContactsByUser(c.Request.Context(), userID.(int64), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"contacts": contacts,
		"count":    len(contacts),
	})
}

// SendMessageRequest is the expected body for sending a manual reply.
type SendMessageRequest struct {
	Content     string `json:"content" binding:"required"`
	MessageType string `json:"message_type"` // defaults to "text"
}

// SendMessage sends a manual reply to a contact via the Meta API.
func (h *InboxHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("user_id")

	contactID, err := strconv.ParseInt(c.Param("contact_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.MessageType == "" {
		req.MessageType = "text"
	}

	ctx := c.Request.Context()

	// Fetch the contact
	contact, err := h.Store.GetContactByID(ctx, contactID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// Verify the contact belongs to this user
	if contact.UserID != userID.(int64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Fetch the channel to get the access token
	channel, err := h.Store.GetChannelByID(ctx, contact.ChannelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Channel not found for contact"})
		return
	}

	// Send via Meta API
	result, err := h.MetaClient.SendMessage(
		contact.Platform,
		channel.AccountID,
		contact.PlatformUserID,
		req.Content,
		channel.AccessToken,
	)
	if err != nil {
		log.Printf("[Inbox] Failed to send message to contact #%d: %v", contactID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	if !result.Success {
		c.JSON(http.StatusBadGateway, gin.H{"error": result.Error})
		return
	}

	// Store the outbound message in DB
	msg := &models.Message{
		UserID:        userID.(int64),
		ChannelID:     channel.ID,
		ContactID:     contact.ID,
		Platform:      contact.Platform,
		Direction:     "outbound",
		Content:       req.Content,
		MessageType:   req.MessageType,
		PlatformMsgID: result.MessageID,
		Status:        "sent",
		IsAutomated:   false,
	}

	if err := h.Store.CreateMessage(ctx, msg); err != nil {
		log.Printf("[Inbox] Message sent but failed to save in DB: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message sent",
		"data": gin.H{
			"id":              msg.ID,
			"content":         msg.Content,
			"platform_msg_id": msg.PlatformMsgID,
			"status":          msg.Status,
			"created_at":      msg.CreatedAt,
		},
	})
}
