package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nguyenanhhao221/greenlight-api/internal/jsonlog"
	"github.com/nguyenanhhao221/greenlight-api/internal/models"
)

// TODO: do this at build time rather than hard code
const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleTime  time.Duration
	}
}
type application struct {
	config config
	logger *jsonlog.Logger
	models models.Models
}

func main() {
	var cfg config

	// Get server config via cli flag
	flag.IntVar(&cfg.port, "port", 42069, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight@localhost/greenlight", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conn", 25, "Max open connection pool for postgres database")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "the duration after which an idle connection will be automatically closed by the health check")

	flag.Parse()

	// Setup our own logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// setup postgres database connection
	logger.Info("Opening database connection using pgxpool", nil)
	connPool, err := openDBConnPool(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer connPool.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: models.New(connPool), // set up basic model for database access layer
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Info("Starting server", map[string]string{"env": cfg.env, "address": srv.Addr})
	if err := srv.ListenAndServe(); err != nil {
		app.logger.PrintFatal(err, nil)
	}
}

func setupDbConfig(cfg config) (*pgxpool.Config, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = int32(cfg.db.maxOpenConns)
	dbConfig.MaxConnIdleTime = cfg.db.maxIdleTime

	return dbConfig, nil
}

func openDBConnPool(cfg config) (*pgxpool.Pool, error) {
	dbConfig, err := setupDbConfig(cfg)
	if err != nil {
		return nil, err
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbpool.Ping(ctxWithTimeout); err != nil {
		return nil, err
	}
	return dbpool, nil
}
