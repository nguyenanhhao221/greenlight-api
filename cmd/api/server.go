package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	shutdownError := make(chan error)

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// That waits (blocks) until a signal is received and sent to the quit channel
		s := <-quit

		app.logger.Info("Shutting down server", "signal", s.String())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call shutdown with 5s timeout context so that the server has a 5 second window to clean up
		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// Log a message to say that we're waiting for any background goroutines to
		// complete their tasks.
		app.logger.Info("completing background tasks", "addr", srv.Addr)
		// Call Wait() to block until our WaitGroup counter is zero --- essentially
		// blocking until the background goroutines have finished. Then we return nil on
		// the shutdownError channel, to indicate that the shutdown completed without
		// any issues.
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("Starting server", "Address", srv.Addr, "environment", app.config.env, "limiter", app.config.limiter)
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.Info("Stopped server", "Address", srv.Addr, "environment", app.config.env)
	return nil
}
