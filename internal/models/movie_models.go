package models

import (
	"context"
	"database/sql"
	"errors"

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

func (m *MovieModel) Get(id int64) (*data.Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

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
	if errors.Is(err, sql.ErrNoRows) {
		return &movie, ErrRecordNotFound
	}
	return &movie, err
}

func (m *MovieModel) Update(movie *data.Movie) error {
	// Use version to prevent data race condition to update. This can be consider as optimistic locking
	query := `
		UPDATE movies 
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version;
	`
	args := []any{movie.Title, movie.Year, movie.Runtime, movie.Genres, movie.ID, movie.Version}

	if err := m.DB.QueryRow(context.Background(), query, args...).Scan(&movie.Version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (m *MovieModel) Delete(id int64) error {
	query := `DELETE FROM movies WHERE id = $1`
	result, err := m.DB.Exec(context.Background(), query, id)
	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}
	return err
}
