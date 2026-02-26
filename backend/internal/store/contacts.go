package store

import (
	"context"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// CreateContact inserts a new contact (lead).
func (s *Storage) CreateContact(ctx context.Context, c *models.Contact) error {
	query := `
		INSERT INTO contacts (user_id, channel_id, platform, platform_user_id, name, phone, email, budget, preferred_location, purchase_timeline, tags, is_hot_lead, booking_state, bot_paused, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	if c.BookingState == "" {
		c.BookingState = "new"
	}
	return s.DB.QueryRow(ctx, query,
		c.UserID, c.ChannelID, c.Platform, c.PlatformUserID,
		c.Name, c.Phone, c.Email, c.Budget,
		c.PreferredLocation, c.PurchaseTimeline, c.Tags,
		c.IsHotLead, c.BookingState, c.BotPaused, now, now,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

// GetOrCreateContact finds a contact by platform user ID for a given user, or creates a new one.
func (s *Storage) GetOrCreateContact(ctx context.Context, c *models.Contact) error {
	query := `
		SELECT id, name, phone, email, budget, preferred_location, purchase_timeline, tags, is_hot_lead, booking_state, bot_paused, created_at, updated_at
		FROM contacts
		WHERE user_id = $1 AND platform = $2 AND platform_user_id = $3`

	err := s.DB.QueryRow(ctx, query, c.UserID, c.Platform, c.PlatformUserID).Scan(
		&c.ID, &c.Name, &c.Phone, &c.Email, &c.Budget,
		&c.PreferredLocation, &c.PurchaseTimeline, &c.Tags,
		&c.IsHotLead, &c.BookingState, &c.BotPaused, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		// Contact doesn't exist, create it
		return s.CreateContact(ctx, c)
	}
	return nil
}

// GetContactsByUser fetches all contacts for a user with pagination.
func (s *Storage) GetContactsByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Contact, error) {
	query := `
		SELECT id, user_id, channel_id, platform, platform_user_id, name, phone, email,
		       budget, preferred_location, purchase_timeline, tags, is_hot_lead, booking_state, bot_paused, created_at, updated_at
		FROM contacts
		WHERE user_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var c models.Contact
		if err := rows.Scan(
			&c.ID, &c.UserID, &c.ChannelID, &c.Platform, &c.PlatformUserID,
			&c.Name, &c.Phone, &c.Email, &c.Budget,
			&c.PreferredLocation, &c.PurchaseTimeline, &c.Tags,
			&c.IsHotLead, &c.BookingState, &c.BotPaused, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

// UpdateContactLead updates the lead qualification fields of a contact.
func (s *Storage) UpdateContactLead(ctx context.Context, contactID int64, budget, location, timeline, phone string, isHot bool) error {
	query := `
		UPDATE contacts
		SET budget = $2, preferred_location = $3, purchase_timeline = $4, phone = $5, is_hot_lead = $6, updated_at = $7
		WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, contactID, budget, location, timeline, phone, isHot, time.Now())
	return err
}

// GetContactByID fetches a single contact by ID.
func (s *Storage) GetContactByID(ctx context.Context, contactID int64) (*models.Contact, error) {
	c := &models.Contact{}
	query := `
		SELECT id, user_id, channel_id, platform, platform_user_id, name, phone, email,
		       budget, preferred_location, purchase_timeline, tags, is_hot_lead, booking_state, bot_paused, created_at, updated_at
		FROM contacts
		WHERE id = $1`

	err := s.DB.QueryRow(ctx, query, contactID).Scan(
		&c.ID, &c.UserID, &c.ChannelID, &c.Platform, &c.PlatformUserID,
		&c.Name, &c.Phone, &c.Email, &c.Budget,
		&c.PreferredLocation, &c.PurchaseTimeline, &c.Tags,
		&c.IsHotLead, &c.BookingState, &c.BotPaused, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// UpdateContactState updates the booking state and automation pause state of a contact.
func (s *Storage) UpdateContactState(ctx context.Context, contactID int64, bookingState string, botPaused bool) error {
	query := `
		UPDATE contacts
		SET booking_state = $2, bot_paused = $3, updated_at = $4
		WHERE id = $1`

	_, err := s.DB.Exec(ctx, query, contactID, bookingState, botPaused, time.Now())
	return err
}
