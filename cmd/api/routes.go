package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// App custom error handler for method not allow and not found
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Map the routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)

	// Movies routes
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)

	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHanlder)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)

	// Users routes
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activate", app.activateUserHandler)

	return app.recoverPanic(app.rateLimitMiddleware(router))
}
