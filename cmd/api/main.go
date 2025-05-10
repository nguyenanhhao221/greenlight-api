package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// TODO: do this at build time rather than hard code
const version = "1.0.0"

type config struct {
	port int
	env  string
}
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	// Get server config via cli flag
	flag.IntVar(&cfg.port, "port", 42069, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Setup our own logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s on port: %d", cfg.env, cfg.port)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
