package store

import (
	"context"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// CreateMessage inserts a new message record.
func (s *Storage) CreateMessage(ctx context.Context, m *models.Message) error {
	query := `
		INSERT INTO messages (user_id, channel_id, contact_id, platform, direction, content, message_type, platform_msg_id, status, is_automated, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`

	return s.DB.QueryRow(ctx, query,
		m.UserID, m.ChannelID, m.ContactID, m.Platform,
		m.Direction, m.Content, m.MessageType, m.PlatformMsgID,
		m.Status, m.IsAutomated, time.Now(),
	).Scan(&m.ID, &m.CreatedAt)
}

// GetMessagesByContact returns messages for a specific contact, ordered by time.
func (s *Storage) GetMessagesByContact(ctx context.Context, contactID int64, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT id, user_id, channel_id, contact_id, platform, direction, content, message_type, platform_msg_id, status, is_automated, created_at
		FROM messages
		WHERE contact_id = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3`

	rows, err := s.DB.Query(ctx, query, contactID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID, &m.UserID, &m.ChannelID, &m.ContactID, &m.Platform,
			&m.Direction, &m.Content, &m.MessageType, &m.PlatformMsgID,
			&m.Status, &m.IsAutomated, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// GetConversations returns the latest message per contact for a user (inbox view).
func (s *Storage) GetConversations(ctx context.Context, userID int64, limit, offset int) ([]models.Message, error) {
	query := `
		SELECT DISTINCT ON (contact_id)
		       id, user_id, channel_id, contact_id, platform, direction, content, message_type, platform_msg_id, status, is_automated, created_at
		FROM messages
		WHERE user_id = $1
		ORDER BY contact_id, created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID, &m.UserID, &m.ChannelID, &m.ContactID, &m.Platform,
			&m.Direction, &m.Content, &m.MessageType, &m.PlatformMsgID,
			&m.Status, &m.IsAutomated, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
