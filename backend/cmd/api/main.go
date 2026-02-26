package main

// @title Lead Automation API
// @version 0.4.0
// @description High-performance Go backend for Lead Automation SaaS.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/social-media-lead/backend/internal/api"
	"github.com/social-media-lead/backend/internal/cache"
	"github.com/social-media-lead/backend/internal/config"
	"github.com/social-media-lead/backend/internal/store"
	"github.com/social-media-lead/backend/internal/workers"
)

func main() {
	// Load .env file (optional, mainly for local dev)
	// Try loading from current dir, then parent, then grandparent (project root)
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			if err := godotenv.Load("../../.env"); err != nil {
				log.Println("No .env file found in ., .., or ../.., reading from environment")
			}
		}
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

	var asynqClient *asynq.Client
	var asynqServer *asynq.Server
	if redisClient != nil {
		redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
		asynqClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, Password: cfg.Redis.Password})
		defer asynqClient.Close()
	}

	// Setup the Gin server
	router, graphWalker := api.SetupRouter(cfg, storage, redisClient, asynqClient)

	if asynqClient != nil {
		redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
		asynqServer = workers.StartServer(redisAddr, graphWalker)
		defer asynqServer.Stop()
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Run server in a goroutine so that it doesn't block the graceful shutdown handling below
	go func() {
		log.Printf("🌐 API server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("❌ Server forced to shutdown: ", err)
	}

	log.Println("✅ Server exiting")
}
