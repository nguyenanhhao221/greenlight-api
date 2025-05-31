package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nguyenanhhao221/greenlight-api/internal/data"
)

type UserModel struct {
	DB *pgxpool.Pool
}

// Create inserts a new user record into the database
func (m UserModel) Create(user *data.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version;
	`
	args := []any{user.Name, user.Email, user.Password.Hash, user.Activated}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRow(ctxWithTimeout, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "users_email_key" {
					return fmt.Errorf("%w: %s", ErrDuplicateEmail, pgErr.Detail)
				}
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByEmail retrieves a user record by email address
func (m UserModel) GetByEmail(email string) (*data.User, error) {
	query := `
	SELECT id, created_at, name, email, password_hash, activated, version 
	FROM users
	WHERE email = $1
	`
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user data.User

	err := m.DB.QueryRow(ctxWithTimeout, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
	}
	return &user, nil
}
