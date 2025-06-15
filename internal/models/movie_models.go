package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nguyenanhhao221/greenlight-api/internal/data"
)

type MovieModel struct {
	DB *pgxpool.Pool
}

func (m MovieModel) Create(movie *data.Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version;
	`
	args := []any{movie.Title, movie.Year, movie.Runtime, movie.Genres}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctxWithTimeout, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) GetAll(title string, genres []string, filters data.Filters) ([]data.Movie, data.Metadata, error) {
	// Use COUNT(*) OVER() to get total count along with the result rows
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER() AS count, id, created_at, title, year, runtime, genres, version	
		FROM movies
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) or $1 = '')
			AND (genres @> $2 or $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4;`,
		filters.SortColumn(), filters.SortDirection(),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	totalRecords := 0
	movie := data.Movie{}
	movies := make([]data.Movie, 0)
	rows, err := m.DB.Query(ctx, query, title, genres, filters.Limit(), filters.Offset())
	if err != nil {
		return nil, data.Metadata{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&totalRecords, &movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, &movie.Genres, &movie.Version)
		if err != nil {
			return nil, data.Metadata{}, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	return movies, data.CalculateMetadata(totalRecords, filters.Page, filters.PageSize), nil
}

func (m MovieModel) Get(id int64) (*data.Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, created_at, title, year, runtime, genres, version FROM movies
	WHERE id=$1;
	`
	movie := data.Movie{}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRow(ctxWithTimeout, query, id).Scan(
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

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := m.DB.QueryRow(ctxWithTimeout, query, args...).Scan(&movie.Version); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (m MovieModel) Delete(id int64) error {
	query := `DELETE FROM movies WHERE id = $1`

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.Exec(ctxWithTimeout, query, id)
	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}
	return err
}
