package main

import (
	"encoding/hex"
	"flag"
	"fmt"
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

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatal(err)
	}

	var cfg config

	// Basic config
	flag.IntVar(&cfg.port, "port", port, "API server port")
	flag.BoolVar(&cfg.debug, "debug", debug, "Debug (true|false)")
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
	emailPortStr := os.Getenv("EMAIL_SERVER_PORT")
	emailPort, err := strconv.Atoi(emailPortStr)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("EMAIL_HOST_SERVER"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", emailPort, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("EMAIL_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("EMAIL_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "GoAuthBackend <no-reply@goauthbackend.johnowolabiidogun.dev>", "SMTP sender")

	// Redis config
	flag.StringVar(&cfg.redisURL, "redis-url", os.Getenv("REDIS_URL"), "Redis URL")

	// Frontend URL
	flag.StringVar(&cfg.frontendURL, "frontend-url", os.Getenv("FRONTEND_URL"), "Frontend URL")
	// Secret
	flag.StringVar(&cfg.secret.HMC, "secret-key", os.Getenv("HMC_SECRET_KEY"), "HMC Secret Key")

	// AWS configs
	flag.StringVar(&cfg.awsConfig.AccessKeyID, "aws-access-key", os.Getenv("AWS_ACCESS_KEY_ID"), "AWS Access KeyID")
	flag.StringVar(&cfg.awsConfig.AccessKeySecret, "aws-access-secret", os.Getenv("AWS_SECRET_ACCESS_KEY"), "AWS Access Secret")
	flag.StringVar(&cfg.awsConfig.Region, "aws-region", os.Getenv("AWS_REGION"), "AWS region")
	flag.StringVar(&cfg.awsConfig.BucketName, "aws-bucketname", os.Getenv("AWS_S3_BUCKET_NAME"), "AWS bucket name")
	flag.Parse()

	cfg.awsConfig.BaseURL = fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com",
		cfg.awsConfig.BucketName,
		cfg.awsConfig.Region,
	)

	cfg.awsConfig.s3_key_prefix = "media/go-auth/"

	secretKey, err := hex.DecodeString(cfg.secret.HMC)
	if err != nil {
		return nil, err
	}
	cfg.secret.secretKey = secretKey
	sessionDuration, err := time.ParseDuration(os.Getenv("SESSION_EXPIRATION"))
	if err != nil {
		return nil, err
	}
	cfg.secret.sessionExpiration = sessionDuration

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
