package main

import (
	"encoding/gob"
	"expvar"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"goauthbackend.johnowolabiidogun.dev/internal/data"
	"goauthbackend.johnowolabiidogun.dev/internal/jsonlog"
	"goauthbackend.johnowolabiidogun.dev/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port  int
	debug bool
	db    struct {
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
		HMC               string
		secretKey         []byte
		sessionExpiration time.Duration
	}
	awsConfig struct {
		AccessKeyID     string
		AccessKeySecret string
		Region          string
		BucketName      string
		BaseURL         string
		s3_key_prefix   string
	}
}

type application struct {
	config      config
	logger      *jsonlog.Logger
	models      data.Models
	mailer      mailer.Mailer
	wg          sync.WaitGroup
	redisClient *redis.Client
	S3Client    *s3.Client
}

func main() {
	gob.Register(&data.UserID{})

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := updateConfigWithEnvVariables()
	if err != nil {
		logger.PrintFatal(err, nil, cfg.debug)
	}

	db, err := openDB(*cfg)

	if err != nil {
		logger.PrintFatal(err, nil, cfg.debug)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil, cfg.debug)

	opt, err := redis.ParseURL(cfg.redisURL)
	if err != nil {
		logger.PrintFatal(err, nil, cfg.debug)
	}

	sdkConfig := aws.Config{
		Region: cfg.awsConfig.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.awsConfig.AccessKeyID, cfg.awsConfig.AccessKeySecret, "",
		),
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	client := redis.NewClient(opt)

	logger.PrintInfo("redis connection pool established", nil, cfg.debug)

	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config:      *cfg,
		logger:      logger,
		models:      data.NewModels(db),
		mailer:      mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		redisClient: client,
		S3Client:    s3Client,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil, cfg.debug)
	}

}
