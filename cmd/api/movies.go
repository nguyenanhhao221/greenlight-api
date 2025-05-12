package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nguyenanhhao221/greenlight-api/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintln(w, "create a new movie"); err != nil {
		log.Printf("error writing to http ResponseWriter: %v\n", err)
	}

	var movieInputData struct {
		Title   string   `json:"title"`   // Movie title
		Year    int32    `json:"year"`    // Movie release year
		Runtime int32    `json:"runtime"` // Movie run time (in minutes)
		Genres  []string `json:"genres"`  // Slice of genres for the movie (romance, comedy, etc.)
	}

	err := json.NewDecoder(r.Body).Decode(&movieInputData)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
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
