package store

import (
	"context"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// UpsertPropertyVisitConfig saves wizard configuration for a user.
// This is an INSERT ... ON CONFLICT UPDATE so the wizard is idempotent.
func (s *Storage) UpsertPropertyVisitConfig(ctx context.Context, cfg *models.PropertyVisitConfig) error {
	query := `
		INSERT INTO property_visit_configs (user_id, project_name, brochure_url, agent_phone, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, TRUE, $5, $5)
		ON CONFLICT (user_id)
		DO UPDATE SET
			project_name = EXCLUDED.project_name,
			brochure_url = EXCLUDED.brochure_url,
			agent_phone  = EXCLUDED.agent_phone,
			is_active    = TRUE,
			updated_at   = EXCLUDED.updated_at
		RETURNING id, created_at, updated_at
	`
	now := time.Now()
	return s.DB.QueryRow(ctx, query,
		cfg.UserID,
		cfg.ProjectName,
		cfg.BrochureURL,
		cfg.AgentPhone,
		now,
	).Scan(&cfg.ID, &cfg.CreatedAt, &cfg.UpdatedAt)
}

// GetPropertyVisitConfig retrieves the active config for a tenant.
// Returns nil, nil when the user has not configured the wizard yet.
func (s *Storage) GetPropertyVisitConfig(ctx context.Context, userID int64) (*models.PropertyVisitConfig, error) {
	query := `
		SELECT id, user_id, project_name, brochure_url, agent_phone, is_active, created_at, updated_at
		FROM property_visit_configs
		WHERE user_id = $1
		LIMIT 1
	`
	var cfg models.PropertyVisitConfig
	err := s.DB.QueryRow(ctx, query, userID).Scan(
		&cfg.ID, &cfg.UserID, &cfg.ProjectName,
		&cfg.BrochureURL, &cfg.AgentPhone, &cfg.IsActive,
		&cfg.CreatedAt, &cfg.UpdatedAt,
	)
	if err != nil {
		// Return nil config (not an error) when no row found
		return nil, nil
	}
	return &cfg, nil
}
