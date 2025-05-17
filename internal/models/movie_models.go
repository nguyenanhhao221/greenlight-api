package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nguyenanhhao221/greenlight-api/internal/data"
)

type MovieModel struct {
	DB *pgxpool.Pool
}

func (m *MovieModel) Create(movie *data.Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version;
	`
	args := []any{movie.Title, movie.Year, movie.Runtime, movie.Genres}

	return m.DB.QueryRow(context.Background(), query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m *MovieModel) Get(id int64) (data.Movie, error) {
	query := `
	SELECT id, created_at, title, year, runtime, genres, version FROM movies
	WHERE id=$1;
	`
	movie := data.Movie{}

	err := m.DB.QueryRow(context.Background(), query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		&movie.Genres,
		&movie.Version,
	)
	return movie, err
}
