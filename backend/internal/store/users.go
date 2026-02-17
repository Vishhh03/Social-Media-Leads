package store

import (
	"context"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// CreateUser inserts a new user into the database.
func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, full_name, company_name, google_id, avatar_url, plan, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	return s.DB.QueryRow(ctx, query,
		user.Email, user.PasswordHash, user.FullName, user.CompanyName,
		user.GoogleID, user.AvatarURL, user.Plan, true, now, now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// GetUserByEmail fetches a user by email address.
func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, full_name, company_name, google_id, avatar_url, plan, is_active, created_at, updated_at
		FROM users
		WHERE email = $1`

	err := s.DB.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.CompanyName, &user.GoogleID, &user.AvatarURL, &user.Plan, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID fetches a user by their ID.
func (s *Storage) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, full_name, company_name, google_id, avatar_url, plan, is_active, created_at, updated_at
		FROM users
		WHERE id = $1`

	err := s.DB.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.CompanyName, &user.GoogleID, &user.AvatarURL, &user.Plan, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetOrCreateOAuthUser finds a user by Google ID or email, or creates a new one.
func (s *Storage) GetOrCreateOAuthUser(ctx context.Context, oauthUser *models.User) (*models.User, error) {
	// Try to find by Google ID first
	if oauthUser.GoogleID != "" {
		user := &models.User{}
		query := `
			SELECT id, email, password_hash, full_name, company_name, google_id, avatar_url, plan, is_active, created_at, updated_at
			FROM users
			WHERE google_id = $1`
		err := s.DB.QueryRow(ctx, query, oauthUser.GoogleID).Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
			&user.CompanyName, &user.GoogleID, &user.AvatarURL, &user.Plan, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err == nil {
			return user, nil
		}
	}

	// Try to find by email (link existing account)
	existing, err := s.GetUserByEmail(ctx, oauthUser.Email)
	if err == nil {
		// Link Google ID to existing account
		_, err = s.DB.Exec(ctx,
			`UPDATE users SET google_id = $2, avatar_url = $3, updated_at = $4 WHERE id = $1`,
			existing.ID, oauthUser.GoogleID, oauthUser.AvatarURL, time.Now())
		if err == nil {
			existing.GoogleID = oauthUser.GoogleID
			existing.AvatarURL = oauthUser.AvatarURL
		}
		return existing, nil
	}

	// Create new user
	oauthUser.Plan = "starter"
	oauthUser.IsActive = true
	if err := s.CreateUser(ctx, oauthUser); err != nil {
		return nil, err
	}
	return oauthUser, nil
}

// UpdateUserProfile updates a user's name, email, and company.
func (s *Storage) UpdateUserProfile(ctx context.Context, userID int64, fullName, email, companyName string) (*models.User, error) {
	query := `
		UPDATE users
		SET full_name = $2, email = $3, company_name = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, email, password_hash, full_name, company_name, google_id, avatar_url, plan, is_active, created_at, updated_at`

	user := &models.User{}
	err := s.DB.QueryRow(ctx, query, userID, fullName, email, companyName, time.Now()).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.CompanyName, &user.GoogleID, &user.AvatarURL, &user.Plan, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUserPassword updates a user's password hash.
func (s *Storage) UpdateUserPassword(ctx context.Context, userID int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1`
	_, err := s.DB.Exec(ctx, query, userID, passwordHash, time.Now())
	return err
}

// UpdateChannelToken updates the access token and expiry for a channel.
func (s *Storage) UpdateChannelToken(ctx context.Context, channelID int64, accessToken string, expiry time.Time) error {
	query := `UPDATE channels SET access_token = $2, token_expiry = $3, updated_at = $4 WHERE id = $1`
	_, err := s.DB.Exec(ctx, query, channelID, accessToken, expiry, time.Now())
	return err
}
