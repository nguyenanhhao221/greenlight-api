package models

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
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

func (m UserModel) GetUserWithToken(plainToken string, scope string) (*data.User, error) {
	slog.Info("GetUserWithToken", "scope", scope)
	tokenHash := sha256.Sum256([]byte(plainToken))
	query := `
	SELECT users.id, users.name, users.email, users.activated, users.created_at, users.version
	FROM  users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	user := data.User{}
	args := []any{tokenHash[:], scope, time.Now()}
	err := m.DB.QueryRow(ctxWithTimeout, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Activated, &user.CreatedAt, &user.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	slog.Info("Successfully get user with provided token")
	return &user, nil
}

func (m UserModel) Update(user *data.User) error {
	slog.Info("Update user to activate in db")
	// Use version to prevent data race condition to update. This can be consider as optimistic locking
	query := `
		UPDATE users 
		SET name = $1, email = $2, activated = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version;
	`
	args := []any{user.Name, user.Email, user.Activated, user.ID, user.Version}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := m.DB.QueryRow(ctxWithTimeout, query, args...).Scan(&user.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return fmt.Errorf("error update user to db: %w", err)
	}
	return nil
}

func (m UserModel) DeleteAllTokenForUser(scope string, userId int64) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1
	AND user_id = $2;
	`
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := m.DB.Exec(ctxWithTimeout, query, scope, userId); err != nil {
		return fmt.Errorf("error DeleteAllTokenForUser %w", err)
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
