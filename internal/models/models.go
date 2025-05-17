package models

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	Movie MovieModel
}

func New(db *pgxpool.Pool) Models {
	return Models{
		Movie: MovieModel{DB: db},
	}
}
