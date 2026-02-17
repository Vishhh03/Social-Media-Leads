package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application.
type Config struct {
	AppPort     string
	AppEnv      string
	GinMode     string
	FrontendURL string
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	Meta        MetaConfig
	Google      GoogleOAuthConfig
}

// GoogleOAuthConfig holds Google OAuth2 settings.
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// DatabaseConfig holds database connection settings.
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
}

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	Secret string
}

// MetaConfig holds Meta/Facebook API settings.
type MetaConfig struct {
	AppID           string
	AppSecret       string
	VerifyToken     string
	PageAccessToken string
	WhatsAppToken   string
}

// Load reads configuration from environment variables.
func Load() *Config {
	cfg := &Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		AppEnv:      getEnv("APP_ENV", "development"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "leadbot"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "leadautomation"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", ""),
		},
		Meta: MetaConfig{
			AppID:           getEnv("META_APP_ID", ""),
			AppSecret:       getEnv("META_APP_SECRET", ""),
			VerifyToken:     getEnv("META_VERIFY_TOKEN", ""),
			PageAccessToken: getEnv("META_PAGE_ACCESS_TOKEN", ""),
			WhatsAppToken:   getEnv("META_WHATSAPP_TOKEN", ""),
		},
		Google: GoogleOAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
		},
	}

	// Build DATABASE_URL if not explicitly set
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		cfg.Database.URL = dbURL
	} else {
		cfg.Database.URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Database.User, cfg.Database.Password,
			cfg.Database.Host, cfg.Database.Port,
			cfg.Database.DBName, cfg.Database.SSLMode,
		)
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
