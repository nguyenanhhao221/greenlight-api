package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/nguyenanhhao221/greenlight-api/internal/data"
	"github.com/nguyenanhhao221/greenlight-api/internal/models"
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

	// After created the movie in database, set the header to include a "Location"
	// so that later the client know where to get the newly created movie
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// Write back to client success response
	if err := app.writeJSON(w, http.StatusCreated, envelop{"movie": movie}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHanlder(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Get movie from database
	movie, err := app.models.Movie.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelop{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Get movie from database
	movie, err := app.models.Movie.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// We use pointers here for the input in order to support partial update
	var movieInputData struct {
		Title   *string       `json:"title"`   // Movie title
		Year    *int32        `json:"year"`    // Movie release year
		Runtime *data.Runtime `json:"runtime"` // Movie run time (in minutes)
		Genres  []string      `json:"genres"`  // Slice of genres for the movie (romance, comedy, etc.)
	}

	err = app.readJSON(w, r, &movieInputData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if movieInputData.Title != nil {
		movie.Title = *movieInputData.Title
	}

	if movieInputData.Year != nil {
		movie.Year = *movieInputData.Year
	}

	if movieInputData.Runtime != nil {
		movie.Runtime = *movieInputData.Runtime
	}

	if movieInputData.Genres != nil {
		movie.Genres = movieInputData.Genres
	}
	// Validate user input
	validator := validator.New()
	if data.ValidateMovie(validator, movie); !validator.Valid() {
		app.failValidationResponse(w, r, validator.Errors)
		return
	}

	// Update movie to database
	err = app.models.Movie.Update(movie)
	if err != nil {
		if errors.Is(err, models.ErrEditConflict) {
			app.editConflictResponse(w,r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelop{"movie": movie}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Get movie from database
	_, err = app.models.Movie.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Movie.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
	}
	if err := app.writeJSON(w, http.StatusOK, envelop{"message": "movie successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
