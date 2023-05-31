package main

import (
	"encoding/hex"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func updateConfigWithEnvVariables() (*config, error) {
	// Load environment variables from `.env` file
	err := godotenv.Load(".env", ".env.development")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	maxOpenConns, err := strconv.Atoi(maxOpenConnsStr)
	if err != nil {
		log.Fatal(err)
	}
	maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS")
	maxIdleConns, err := strconv.Atoi(maxIdleConnsStr)
	if err != nil {
		log.Fatal(err)
	}

	emailUsername := os.Getenv("EMAIL_USERNAME")
	emailPassword := os.Getenv("EMAIL_PASSWORD")

	var cfg config

	// Basic config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// Database config
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", maxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", maxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime,
		"db-max-idle-time",
		os.Getenv("DB_MAX_IDLE_TIME"),
		"PostgreSQL max connection idle time",
	)
	// Email
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", emailUsername, "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", emailPassword, "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "GoAuthBackend <no-reply@goauthbackend.johnowolabiidogun.dev>", "SMTP sender")

	// Redis config
	flag.StringVar(&cfg.redisURL, "redis-url", os.Getenv("REDIS_URL"), "Redis URL")

	// Frontend URL
	flag.StringVar(&cfg.frontendURL, "frontend-url", os.Getenv("FRONTEND_URL"), "Frontend URL")
	// Secret
	hmc_secret := os.Getenv("HMC_SECRET_KEY")
	flag.StringVar(&cfg.secret.HMC, "secret-key", hmc_secret, "HMC Secret Key")
	flag.Parse()

	secretKey, err := hex.DecodeString(cfg.secret.HMC)
	if err != nil {
		return nil, err
	}
	cfg.secret.secretKey = secretKey

	// Token Expiration
	tokexpirationStr := os.Getenv("TOKEN_EXPIRATION")
	duration, err := time.ParseDuration(tokexpirationStr)
	if err != nil {
		return nil, err
	}
	cfg.tokenExpiration.durationString = tokexpirationStr
	cfg.tokenExpiration.duration = duration

	return &cfg, nil
}
