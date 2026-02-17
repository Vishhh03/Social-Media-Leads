package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Storage wraps the database connection pool.
type Storage struct {
	DB *pgxpool.Pool
}

// New creates a new Storage with a connection pool.
func New(databaseURL string) (*Storage, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Storage{DB: pool}, nil
}

// Close closes the database connection pool.
func (s *Storage) Close() {
	s.DB.Close()
}

// RunMigrations executes the SQL migration files.
func (s *Storage) RunMigrations() error {
	migrationSQL, err := os.ReadFile("internal/store/migrations/001_initial_schema.sql")
	if err != nil {
		return fmt.Errorf("unable to read migration file: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = s.DB.Exec(ctx, string(migrationSQL))
	if err != nil {
		return fmt.Errorf("unable to run migration: %w", err)
	}

	return nil
}
