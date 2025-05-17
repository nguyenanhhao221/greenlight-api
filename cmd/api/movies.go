package main

import (
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

	// Validate user input
	validator := validator.New()

	movie := data.Movie{
		Title:   movieInputData.Title,
		Year:    movieInputData.Year,
		Runtime: movieInputData.Runtime,
		Genres:  movieInputData.Genres,
	}
	if data.ValidateMovie(validator, &movie); !validator.Valid() {
		app.failValidationResponse(w, r, validator.Errors)
		return
	}

	// Create movie to database, Create will also update the fill the movie variable with correct data
	err = app.models.Movie.Create(&movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Write back to client success response
	if err := app.writeJSON(w, http.StatusCreated, envelop{"movie": movie}, nil); err != nil {
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
