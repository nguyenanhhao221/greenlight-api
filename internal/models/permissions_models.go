package models

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionModel struct {
	DB *pgxpool.Pool
}

type Permissions []string

func (m PermissionModel) GetAllForUser(userId int64) (Permissions, error) {
	query := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
	INNER JOIN users ON users_permissions.user_id = users.id
	WHERE users.id = $1;
	`

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.Query(ctxWithTimeout, query, userId)
	if err != nil {
		err = fmt.Errorf("error when Query in GetAllForUser %w", err)
		return nil, err
	}
	permissions, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		err = fmt.Errorf("pgx.CollectRows error %w", err)
		return nil, err
	}
	return permissions, nil
}

// Includes is a helper to check if permission provide exist in the user's permissions
func (p Permissions) Includes(permission string) bool {
	return slices.Contains(p, permission)
}
