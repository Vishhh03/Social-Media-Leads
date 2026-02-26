package store

import (
	"context"

	"github.com/social-media-lead/backend/internal/models"
)

// CreateVisit inserts a new visit booking into the database
func (s *Storage) CreateVisit(ctx context.Context, v *models.Visit) error {
	query := `
		INSERT INTO visits (user_id, contact_id, project_name, visit_time, status, lead_source_channel)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return s.DB.QueryRow(ctx, query,
		v.UserID,
		v.ContactID,
		v.ProjectName,
		v.VisitTime,
		v.Status,
		v.LeadSourceChannel,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
}

func (s *Storage) GetVisitsByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Visit, error) {
	query := `
		SELECT id, user_id, contact_id, project_name, visit_time, status, lead_source_channel, created_at, updated_at
		FROM visits
		WHERE user_id = $1
		ORDER BY visit_time ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := s.DB.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var visits []models.Visit
	for rows.Next() {
		var v models.Visit
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.ContactID, &v.ProjectName,
			&v.VisitTime, &v.Status, &v.LeadSourceChannel,
			&v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, err
		}
		visits = append(visits, v)
	}
	return visits, rows.Err()
}

func (s *Storage) GetVisitByContact(ctx context.Context, contactID int64) (*models.Visit, error) {
	query := `
		SELECT id, user_id, contact_id, project_name, visit_time, status, lead_source_channel, created_at, updated_at
		FROM visits
		WHERE contact_id = $1
		ORDER BY created_at DESC LIMIT 1
	`
	var v models.Visit
	err := s.DB.QueryRow(ctx, query, contactID).Scan(
		&v.ID, &v.UserID, &v.ContactID, &v.ProjectName,
		&v.VisitTime, &v.Status, &v.LeadSourceChannel,
		&v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (s *Storage) UpdateVisitStatus(ctx context.Context, visitID int64, status string) error {
	query := `
		UPDATE visits
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := s.DB.Exec(ctx, query, status, visitID)
	return err
}
