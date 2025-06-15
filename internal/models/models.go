package models

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	Movie      MovieModel
	User       UserModel
	Token      TokenModel
	Permission PermissionModel
}

func New(db *pgxpool.Pool) Models {
	return Models{
		Movie:      MovieModel{DB: db},
		User:       UserModel{DB: db},
		Token:      TokenModel{DB: db},
		Permission: PermissionModel{DB: db},
	}
}
