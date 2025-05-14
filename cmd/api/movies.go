package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nguyenanhhao221/greenlight-api/internal/data"
	"github.com/nguyenanhhao221/greenlight-api/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var movieInputData struct {
		Title   string       `json:"title"`   // Movie title
		Year    int32        `json:"year"`    // Movie release year
		Runtime data.Runtime `json:"runtime"` // Movie run time (in minutes)
		Genres  []string     `json:"genres"`  // Slice of genres for the movie (romance, comedy, etc.)
	}

	err := app.readJSON(w, r, &movieInputData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	validator := validator.New()
	validator.Check(movieInputData.Title != "", "title", "must be provided")
	validator.Check(len(movieInputData.Title) <= 500, "title", "must not be 500 bytes long")

	validator.Check(movieInputData.Year != 0, "year", "must be provided")
	validator.Check(movieInputData.Year >= 1888, "year", "must be greater than 1888")
	validator.Check(movieInputData.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	validator.Check(movieInputData.Runtime != 0, "runtime", "must be provided")
	validator.Check(movieInputData.Runtime > 0, "runtime", "must be positive integer")

	validator.Check(movieInputData.Genres != nil, "genres", "must be provided")
	validator.Check(len(movieInputData.Genres) >= 1, "genres", "must contains at least 1 genre")
	validator.Check(len(movieInputData.Genres) <= 5, "genres", "must not contains more than 5 genres")
	validator.Check(len(movieInputData.Genres) <= 5, "genres", "must not contains more than 5 genres")
	validator.Check(validator.Unique(movieInputData.Genres), "genres", "must contains unique values")

	if !validator.Valid() {
		app.failValidationResponse(w, r, validator.Errors)
		return
	}

	_, err = fmt.Fprintf(w, "%+v\n", movieInputData)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHanlder(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movieData := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Batman",
		Runtime:   102,
		Genres:    []string{"superhero", "fiction"},
		Version:   1,
	}

	if err := app.writeJSON(w, http.StatusOK, envelop{"movie": movieData}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
