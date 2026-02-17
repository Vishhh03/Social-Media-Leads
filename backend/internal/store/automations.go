package store

import (
	"context"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// CreateAutomation inserts a new automation rule.
func (s *Storage) CreateAutomation(ctx context.Context, a *models.Automation) error {
	query := `
		INSERT INTO automations (user_id, name, trigger_type, keywords, reply_text, reply_media, delay_ms, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	return s.DB.QueryRow(ctx, query,
		a.UserID, a.Name, a.TriggerType, a.Keywords,
		a.ReplyText, a.ReplyMedia, a.DelayMs, true, now, now,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

// GetAutomationsByUser fetches all active automations for a user.
func (s *Storage) GetAutomationsByUser(ctx context.Context, userID int64) ([]models.Automation, error) {
	query := `
		SELECT id, user_id, name, trigger_type, keywords, reply_text, reply_media, delay_ms, is_active, created_at, updated_at
		FROM automations
		WHERE user_id = $1 AND is_active = TRUE
		ORDER BY created_at DESC`

	rows, err := s.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var automations []models.Automation
	for rows.Next() {
		var a models.Automation
		if err := rows.Scan(
			&a.ID, &a.UserID, &a.Name, &a.TriggerType, &a.Keywords,
			&a.ReplyText, &a.ReplyMedia, &a.DelayMs, &a.IsActive,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		automations = append(automations, a)
	}
	return automations, nil
}

// UpdateAutomation updates an existing automation rule.
func (s *Storage) UpdateAutomation(ctx context.Context, a *models.Automation) error {
	query := `
		UPDATE automations
		SET name = $2, trigger_type = $3, keywords = $4, reply_text = $5, reply_media = $6, delay_ms = $7, is_active = $8, updated_at = $9
		WHERE id = $1 AND user_id = $10`

	_, err := s.DB.Exec(ctx, query,
		a.ID, a.Name, a.TriggerType, a.Keywords,
		a.ReplyText, a.ReplyMedia, a.DelayMs, a.IsActive,
		time.Now(), a.UserID,
	)
	return err
}

// DeleteAutomation soft-deletes an automation by deactivating it.
func (s *Storage) DeleteAutomation(ctx context.Context, automationID, userID int64) error {
	query := `UPDATE automations SET is_active = FALSE, updated_at = $3 WHERE id = $1 AND user_id = $2`
	_, err := s.DB.Exec(ctx, query, automationID, userID, time.Now())
	return err
}
