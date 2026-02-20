package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

// AutomationHandler handles automation rule endpoints.
type AutomationHandler struct {
	Store store.Store
}

// CreateAutomationRequest is the expected body for creating an automation.
type CreateAutomationRequest struct {
	Name        string   `json:"name" binding:"required"`
	TriggerType string   `json:"trigger_type" binding:"required"`
	Keywords    []string `json:"keywords"`
	ReplyText   string   `json:"reply_text" binding:"required"`
	ReplyMedia  string   `json:"reply_media"`
	DelayMs     int      `json:"delay_ms"`
}

// ListAutomations returns all active automations for the current user.
func (h *AutomationHandler) ListAutomations(c *gin.Context) {
	userID, _ := c.Get("user_id")

	automations, err := h.Store.GetAutomationsByUser(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch automations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"automations": automations,
		"count":       len(automations),
	})
}

// CreateAutomation creates a new automation rule.
func (h *AutomationHandler) CreateAutomation(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateAutomationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	automation := &models.Automation{
		UserID:      userID.(int64),
		Name:        req.Name,
		TriggerType: req.TriggerType,
		Keywords:    req.Keywords,
		ReplyText:   req.ReplyText,
		ReplyMedia:  req.ReplyMedia,
		DelayMs:     req.DelayMs,
		IsActive:    true,
	}

	if err := h.Store.CreateAutomation(c.Request.Context(), automation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create automation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Automation created",
		"automation": automation,
	})
}

// DeleteAutomation deactivates an automation rule.
func (h *AutomationHandler) DeleteAutomation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	automationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid automation ID"})
		return
	}

	if err := h.Store.DeleteAutomation(c.Request.Context(), automationID, userID.(int64)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete automation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Automation deleted"})
}
