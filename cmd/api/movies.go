package main

import (
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
}

func (app *application) showMovieHanlder(w http.ResponseWriter, r *http.Request) {
	id, err := app.readParams(r)
	if err != nil {
		http.NotFound(w, r)
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
		app.logger.Printf("error writing to http ResponseWriter: %v\n", err)
		http.Error(w, "The server encounter a problem and couldn't process your request", http.StatusInternalServerError)
	}
}
