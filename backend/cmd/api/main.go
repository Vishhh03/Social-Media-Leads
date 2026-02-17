package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/social-media-lead/backend/internal/api"
	"github.com/social-media-lead/backend/internal/cache"
	"github.com/social-media-lead/backend/internal/config"
	"github.com/social-media-lead/backend/internal/store"
)

func main() {
	// Load .env file (optional, mainly for local dev)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	// Load configuration
	cfg := config.Load()

	log.Printf("🚀 Starting Lead Automation API (env: %s)", cfg.AppEnv)

	// Connect to database
	storage, err := store.New(cfg.Database.URL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer storage.Close()
	log.Println("✅ Connected to PostgreSQL")

	// Run migrations
	if err := storage.RunMigrations(); err != nil {
		log.Printf("⚠️  Migration warning: %v", err)
	} else {
		log.Println("✅ Database migrations applied")
	}

	// Connect to Redis
	var redisClient *cache.RedisClient
	redisClient, err = cache.New(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Printf("⚠️  Redis connection failed (non-fatal): %v", err)
		redisClient = nil
	} else {
		defer redisClient.Close()
		log.Println("✅ Connected to Redis")
	}

	// Setup and start the Gin server
	router := api.SetupRouter(cfg, storage, redisClient)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🌐 API server listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
