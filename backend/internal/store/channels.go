package store

import (
	"context"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// CreateChannel inserts a new channel for a user.
func (s *Storage) CreateChannel(ctx context.Context, ch *models.Channel) error {
	query := `
		INSERT INTO channels (user_id, platform, account_id, account_name, access_token, refresh_token, token_expiry, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	return s.DB.QueryRow(ctx, query,
		ch.UserID, ch.Platform, ch.AccountID, ch.AccountName,
		ch.AccessToken, ch.RefreshToken, ch.TokenExpiry,
		true, now, now,
	).Scan(&ch.ID, &ch.CreatedAt, &ch.UpdatedAt)
}

// GetChannelsByUser fetches all channels for a given user.
func (s *Storage) GetChannelsByUser(ctx context.Context, userID int64) ([]models.Channel, error) {
	query := `
		SELECT id, user_id, platform, account_id, account_name, is_active, created_at, updated_at
		FROM channels
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := s.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []models.Channel
	for rows.Next() {
		var ch models.Channel
		if err := rows.Scan(
			&ch.ID, &ch.UserID, &ch.Platform, &ch.AccountID,
			&ch.AccountName, &ch.IsActive, &ch.CreatedAt, &ch.UpdatedAt,
		); err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}
	return channels, nil
}

// GetChannelByAccountID finds a channel by its platform and account_id.
// This is used during webhook processing to resolve which user owns the account.
func (s *Storage) GetChannelByAccountID(ctx context.Context, platform, accountID string) (*models.Channel, error) {
	ch := &models.Channel{}
	query := `
		SELECT id, user_id, platform, account_id, account_name, access_token, refresh_token, is_active, created_at, updated_at
		FROM channels
		WHERE platform = $1 AND account_id = $2 AND is_active = TRUE
		LIMIT 1`

	err := s.DB.QueryRow(ctx, query, platform, accountID).Scan(
		&ch.ID, &ch.UserID, &ch.Platform, &ch.AccountID, &ch.AccountName,
		&ch.AccessToken, &ch.RefreshToken, &ch.IsActive, &ch.CreatedAt, &ch.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// GetChannelByID fetches a channel by its ID.
func (s *Storage) GetChannelByID(ctx context.Context, channelID int64) (*models.Channel, error) {
	ch := &models.Channel{}
	query := `
		SELECT id, user_id, platform, account_id, account_name, access_token, refresh_token, is_active, created_at, updated_at
		FROM channels
		WHERE id = $1`

	err := s.DB.QueryRow(ctx, query, channelID).Scan(
		&ch.ID, &ch.UserID, &ch.Platform, &ch.AccountID, &ch.AccountName,
		&ch.AccessToken, &ch.RefreshToken, &ch.IsActive, &ch.CreatedAt, &ch.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// DeleteChannel soft-deletes a channel by deactivating it.
func (s *Storage) DeleteChannel(ctx context.Context, channelID, userID int64) error {
	query := `UPDATE channels SET is_active = FALSE, updated_at = $3 WHERE id = $1 AND user_id = $2`
	_, err := s.DB.Exec(ctx, query, channelID, userID, time.Now())
	return err
}
