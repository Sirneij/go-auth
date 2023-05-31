package main

import (
	"encoding/gob"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"goauthbackend.johnowolabiidogun.dev/internal/data"
	"goauthbackend.johnowolabiidogun.dev/internal/jsonlog"
	"goauthbackend.johnowolabiidogun.dev/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	redisURL        string
	frontendURL     string
	tokenExpiration struct {
		durationString string
		duration       time.Duration
	}
	secret struct {
		HMC       string
		secretKey []byte
	}
}

type application struct {
	config      config
	logger      *jsonlog.Logger
	models      data.Models
	mailer      mailer.Mailer
	wg          sync.WaitGroup
	redisClient *redis.Client
}

func main() {
	gob.Register(&data.UserID{})

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := updateConfigWithEnvVariables()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := openDB(*cfg)

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	opt, err := redis.ParseURL(cfg.redisURL)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	client := redis.NewClient(opt)

	logger.PrintInfo("redis connection pool established", nil)

	app := &application{
		config:      *cfg,
		logger:      logger,
		models:      data.NewModels(db),
		mailer:      mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		redisClient: client,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}
