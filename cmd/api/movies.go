package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nguyenanhhao221/greenlight-api/internal/data"
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
