package repository

import (
	"feklistova/models"
	"context"
	"time"
)

// CreateSession creates a new session in the database and returns the created session
func (r *Repository) CreateSession(ctx context.Context, userID int, expirationTime time.Time, ipAddress string, userAgent string) (*models.Session, error) {
	var sessionID int
	err := r.Db.QueryRowContext(ctx, `
        INSERT INTO sessions (user_id, expiration_time, ip_address, user_agent, created_at)
        VALUES ($1, $2, $3, $4, NOW())
        RETURNING session_id`,
		userID, expirationTime, ipAddress, userAgent).Scan(&sessionID)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		SessionID:      sessionID,
		UserID:         userID,
		ExpirationTime: expirationTime,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		CreatedAt:      time.Now(),
	}

	return session, nil
}

// DeleteSession deletes a session from the database by its session ID
func (r *Repository) DeleteSession(ctx context.Context, sessionID int) error {
	_, err := r.Db.ExecContext(ctx, "DELETE FROM sessions WHERE session_id = $1", sessionID)
	if err != nil {
		return err
	}
	return nil
}
