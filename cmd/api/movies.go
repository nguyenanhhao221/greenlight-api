package main

import (
	"fmt"
	"log"
	"net/http"
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

	if _, err := fmt.Fprintf(w, "show the details of movies %d\n", id); err != nil {
		log.Printf("error writing to http ResponseWriter: %v\n", err)
	}
}
