package store

import (
	"context"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// CreateBroadcast inserts a new broadcast draft.
func (s *Storage) CreateBroadcast(ctx context.Context, b *models.Broadcast) error {
	query := `
		INSERT INTO broadcasts (user_id, name, content, media_url, status, total_sent, total_failed, scheduled_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	return s.DB.QueryRow(ctx, query,
		b.UserID, b.Name, b.Content, b.MediaURL,
		"draft", 0, 0, b.ScheduledAt, now, now,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

// GetBroadcastsByUser fetches all broadcasts for a user.
func (s *Storage) GetBroadcastsByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Broadcast, error) {
	query := `
		SELECT id, user_id, name, content, media_url, status, total_sent, total_failed, scheduled_at, sent_at, created_at, updated_at
		FROM broadcasts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var broadcasts []models.Broadcast
	for rows.Next() {
		var b models.Broadcast
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.Name, &b.Content, &b.MediaURL,
			&b.Status, &b.TotalSent, &b.TotalFailed,
			&b.ScheduledAt, &b.SentAt, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		broadcasts = append(broadcasts, b)
	}
	return broadcasts, nil
}

// GetBroadcastByID fetches a single broadcast by ID.
func (s *Storage) GetBroadcastByID(ctx context.Context, broadcastID int64) (*models.Broadcast, error) {
	b := &models.Broadcast{}
	query := `
		SELECT id, user_id, name, content, media_url, status, total_sent, total_failed, scheduled_at, sent_at, created_at, updated_at
		FROM broadcasts
		WHERE id = $1`

	err := s.DB.QueryRow(ctx, query, broadcastID).Scan(
		&b.ID, &b.UserID, &b.Name, &b.Content, &b.MediaURL,
		&b.Status, &b.TotalSent, &b.TotalFailed,
		&b.ScheduledAt, &b.SentAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateBroadcastStatus updates the status and counters of a broadcast.
func (s *Storage) UpdateBroadcastStatus(ctx context.Context, broadcastID int64, status string, totalSent, totalFailed int) error {
	query := `
		UPDATE broadcasts
		SET status = $2, total_sent = $3, total_failed = $4, sent_at = $5, updated_at = $6
		WHERE id = $1`

	now := time.Now()
	_, err := s.DB.Exec(ctx, query, broadcastID, status, totalSent, totalFailed, now, now)
	return err
}
