package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	js := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	if err := app.writeJSON(w, http.StatusOK, js, nil); err != nil {
		app.logger.Printf("error writing to JSON: %v\n", err)
		http.Error(w, "The server encounter a problem and couldn't process your request", http.StatusInternalServerError)
	}
}
