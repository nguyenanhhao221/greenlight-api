package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/nguyenanhhao221/greenlight-api/internal/mailer"
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
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}
type application struct {
	config config
	logger *slog.Logger
	models models.Models
	mailer *mailer.Mailer
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
		os.Exit(1)
	}
	mailPassword := os.Getenv("MAIL_SERVER_PASSWORD")
	if mailPassword == "" {
		slog.Error("need to set environment variable: MAIL_SERVER_PASSWORD")
		os.Exit(1)
	}

	var cfg config

	// Get server config via cli flag
	flag.IntVar(&cfg.port, "port", 42069, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight@localhost/greenlight", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conn", 25, "Max open connection pool for postgres database")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "the duration after which an idle connection will be automatically closed by the health check")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	// Setup for smtp configuration, credential need to be set up via MailTrap
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "633e69db2a27b2", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", mailPassword, "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "hao@haonguyen.tech", "SMTP sender")
	flag.Parse()

	// Initialize default slog
	slogger := slog.Default()

	// setup postgres database connection
	slog.Info("Opening database connection using pgxpool")
	connPool, err := openDBConnPool(cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer connPool.Close()

	mailer, err := mailer.New(cfg.smtp.host, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)
	if err != nil {
		slog.Error("error setting up mailer: ", "err:", err.Error())
		os.Exit(1)
	}
	app := &application{
		config: cfg,
		logger: slogger,
		models: models.New(connPool), // set up basic model for database access layer
		mailer: mailer,
	}

	if err := app.serve(); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
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
