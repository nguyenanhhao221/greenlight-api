package models

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound = errors.New("record not found") 
	ErrEditConflict = errors.New("edit conflict") 
)

type Models struct {
	Movie MovieModel
}

func New(db *pgxpool.Pool) Models {
	return Models{
		Movie: MovieModel{DB: db},
	}
}
