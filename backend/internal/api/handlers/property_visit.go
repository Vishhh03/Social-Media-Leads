package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/cache"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

// PropertyVisitHandler manages the wizard configuration for property visits.
type PropertyVisitHandler struct {
	Store store.Store
	Cache *cache.RedisClient
}

// ActivateRequest is the body the wizard POSTs on activation.
type ActivateRequest struct {
	ProjectName string `json:"project_name" binding:"required"`
	BrochureURL string `json:"brochure_url"`
	AgentPhone  string `json:"agent_phone" binding:"required"`
}

// Activate saves the wizard config and invalidates the Redis cache so the
// next inbound webhook picks up the fresh values immediately.
func (h *PropertyVisitHandler) Activate(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := &models.PropertyVisitConfig{
		UserID:      userID.(int64),
		ProjectName: req.ProjectName,
		BrochureURL: req.BrochureURL,
		AgentPhone:  req.AgentPhone,
		IsActive:    true,
	}

	ctx := c.Request.Context()
	if err := h.Store.UpsertPropertyVisitConfig(ctx, cfg); err != nil {
		log.Printf("[PropertyVisit] Failed to upsert config for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save configuration"})
		return
	}

	// Bust the Redis cache so the webhook engine reloads immediately
	if h.Cache != nil {
		if err := h.Cache.InvalidateVisitConfig(ctx, userID.(int64)); err != nil {
			log.Printf("[PropertyVisit] Cache invalidation failed (non-fatal): %v", err)
		}
	}

	log.Printf("[PropertyVisit] ✅ Config activated for user %d: project=%s", userID, cfg.ProjectName)
	c.JSON(http.StatusOK, gin.H{
		"message": "Automation activated",
		"config":  cfg,
	})
}

// GetConfig returns the current wizard configuration for the user.
func (h *PropertyVisitHandler) GetConfig(c *gin.Context) {
	userID, _ := c.Get("user_id")
	ctx := c.Request.Context()

	cfg, err := h.Store.GetPropertyVisitConfig(ctx, userID.(int64))
	if err != nil || cfg == nil {
		c.JSON(http.StatusOK, gin.H{"config": nil, "active": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"config": cfg, "active": cfg.IsActive})
}

