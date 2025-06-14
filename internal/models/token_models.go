package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nguyenanhhao221/greenlight-api/internal/data"
)

type TokenModel struct {
	DB *pgxpool.Pool
}

func (m TokenModel) generateToken() (string, []byte) {
	plainRandomStr := rand.Text()
	hash := sha256.Sum256([]byte(plainRandomStr))
	return plainRandomStr, hash[:]
}

// New Generate an activation token for the newly created user and insert into token table
func (m TokenModel) New(userId int64, ttl time.Duration, scope string) (*data.Token, error) {
	plain, hash := m.generateToken()
	expiry := time.Now().Add(ttl)

	token := &data.Token{UserID: userId, Plain: plain, Hash: hash, Expiry: expiry, Scope: scope}
	if err := m.create(token); err != nil {
		return nil, fmt.Errorf("m.create token error :%w", err)
	}
	return token, nil
}

func (m TokenModel) create(token *data.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4);
	`
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.Exec(ctxWithTimeout, query, args...)
	return err
}

func (m TokenModel) GetByUserId(userId int64) (*data.Token, error) {
	if userId < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT hash, user_id, expiry, scope FROM tokens
	WHERE user_id=$1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	token := data.Token{}
	err := m.DB.QueryRow(ctx, query, userId).Scan(&token.Hash, &token.UserID, &token.Expiry, &token.Scope)
	if errors.Is(err, sql.ErrNoRows) {
		return &token, ErrRecordNotFound
	}
	return &token, err
}
