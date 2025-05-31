package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     slog.NewLogLogger(slog.NewTextHandler(os.Stderr, nil), slog.LevelError),
	}

	app.logger.Info("Starting server", "Address", srv.Addr, "environment", app.config.env, "limiter", app.config.limiter)
	return srv.ListenAndServe()
}
