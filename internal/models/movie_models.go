package models

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieModel struct {
	DB *pgxpool.Pool
}

func (m *MovieModel) Create() {
}
