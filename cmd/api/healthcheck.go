package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	js := `{"status": "available" , "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write([]byte(js)); err != nil {
		app.logger.Printf("error wring to http.ResponseWriter: %v\n", err)
		return
	}
}
