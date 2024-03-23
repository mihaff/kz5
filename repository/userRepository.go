package repository

import (
	"feklistova/models"
	"context"
)

// GetUserByID retrieves a user from the database by ID.
func (r *Repository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	query := "SELECT user_id, username FROM users WHERE user_id = $1"
	row := r.Db.QueryRowContext(ctx, query, id)
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user from the database by ID.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT user_id, username, email, password, created_at FROM users WHERE email = $1"

	user := &models.User{}
	err := r.Db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// RegisterUser creates a new user in the database.
func (r *Repository) RegisterUser(ctx context.Context, user models.User) (int, error) {
	query := `
        INSERT INTO users (username, email, password, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING user_id
    `
	var userID int
	err := r.Db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, user.CreatedAt).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
