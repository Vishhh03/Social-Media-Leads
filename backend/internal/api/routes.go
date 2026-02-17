package api

import (
	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/api/handlers"
	"github.com/social-media-lead/backend/internal/api/middleware"
	"github.com/social-media-lead/backend/internal/cache"
	"github.com/social-media-lead/backend/internal/config"
	"github.com/social-media-lead/backend/internal/meta"
	"github.com/social-media-lead/backend/internal/store"
)

// SetupRouter creates and configures the Gin engine with all routes.
func SetupRouter(cfg *config.Config, storage *store.Storage, redisClient *cache.RedisClient) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		redisOk := redisClient != nil
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "lead-automation-api",
			"version": "0.3.0",
			"redis":   redisOk,
		})
	})

	// Initialize Meta API client + token refresher
	metaClient := meta.NewClient()
	tokenRefresher := meta.NewTokenRefresher(cfg.Meta.AppID, cfg.Meta.AppSecret, redisClient)

	// Initialize handlers
	authHandler := &handlers.AuthHandler{Store: storage, JWTSecret: cfg.JWT.Secret}
	webhookHandler := &handlers.WebhookHandler{Store: storage, Config: cfg, MetaClient: metaClient}
	inboxHandler := &handlers.InboxHandler{Store: storage, MetaClient: metaClient}
	automationHandler := &handlers.AutomationHandler{Store: storage}
	channelHandler := &handlers.ChannelHandler{Store: storage, TokenRefresher: tokenRefresher}
	broadcastHandler := &handlers.BroadcastHandler{Store: storage, MetaClient: metaClient, Redis: redisClient}

	// --- Public Routes ---
	v1 := r.Group("/api/v1")
	{
		// Auth (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
		}

		// Meta Webhooks (public, verified by Meta token)
		webhooks := v1.Group("/webhooks")
		{
			webhooks.GET("/meta", webhookHandler.VerifyWebhook)
			webhooks.POST("/meta", webhookHandler.HandleWebhook)
		}
	}

	// --- Protected Routes (JWT required) ---
	protected := v1.Group("")
	protected.Use(middleware.AuthRequired(cfg.JWT.Secret))
	{
		// User profile
		protected.GET("/me", authHandler.Me)
		protected.PUT("/me", authHandler.UpdateProfile)
		protected.PUT("/me/password", authHandler.ChangePassword)

		// Channels
		channels := protected.Group("/channels")
		{
			channels.POST("", channelHandler.ConnectChannel)
			channels.GET("", channelHandler.ListChannels)
			channels.DELETE("/:id", channelHandler.DisconnectChannel)
		}

		// Inbox
		inbox := protected.Group("/inbox")
		{
			inbox.GET("/conversations", inboxHandler.GetConversations)
			inbox.GET("/messages/:contact_id", inboxHandler.GetMessages)
			inbox.POST("/messages/:contact_id", inboxHandler.SendMessage)
			inbox.GET("/contacts", inboxHandler.GetContacts)
		}

		// Automations
		automations := protected.Group("/automations")
		{
			automations.GET("", automationHandler.ListAutomations)
			automations.POST("", automationHandler.CreateAutomation)
			automations.DELETE("/:id", automationHandler.DeleteAutomation)
		}

		// Broadcasts
		broadcasts := protected.Group("/broadcasts")
		{
			broadcasts.GET("", broadcastHandler.ListBroadcasts)
			broadcasts.POST("", broadcastHandler.CreateBroadcast)
			broadcasts.POST("/:id/send", broadcastHandler.SendBroadcast)
		}
	}

	return r
}
