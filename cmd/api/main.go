package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"expvar"
	"flag"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/olzzhas/qrent/internal/data"
	"github.com/olzzhas/qrent/pkg/jsonlog"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
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
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type application struct {
	config      config
	logger      *jsonlog.Logger
	models      data.Models
	redis       *redis.Client
	wg          sync.WaitGroup
	rabbitMQ    *amqp.Connection
	mongoClient *mongo.Client
}

// @title QRent API
// @version 1.0
// @description API для работы с организациями, повербанками и станциями.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:4000
// @BasePath /v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла: ", err)
	}

	cfg := loadConfig()

	flag.IntVar(&cfg.port, "port", cfg.port, "Порт API сервера")
	flag.StringVar(&cfg.env, "env", cfg.env, "Окружение: development|staging|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", cfg.db.dsn, "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", cfg.db.maxOpenConns, "Максимальное число открытых подключений к PostgreSQL")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", cfg.db.maxIdleConns, "Максимальное число простаивающих подключений к PostgreSQL")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", cfg.db.maxIdleTime, "Время простаивания подключения PostgreSQL")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", cfg.limiter.rps, "Максимальное число запросов в секунду для rate limiter")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", cfg.limiter.burst, "Максимальное burst для rate limiter")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", cfg.limiter.enabled, "Включить rate limiter")
	flag.Parse()

	logger, err := jsonlog.New(os.Stdout, jsonlog.LevelInfo, getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"))
	if err != nil {
		log.Fatal("Ошибка при создании логгера: ", err)
	}

	redisClient, err := redisConnect()
	if err != nil {
		logger.PrintFatal(err, nil, "general")
	}
	logger.PrintInfo("Соединение с Redis установлено", nil, "general")

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil, "general")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.PrintFatal(err, nil, "general")
		}
	}()
	logger.PrintInfo("Пул подключений к БД установлен", nil, "general")

	mongoClient, err := connectMongoDB()
	if err != nil {
		logger.PrintFatal(err, nil, "general")
	}
	logger.PrintInfo("Соединение с MongoDB установлено", nil, "general")

	rabbitConn, err := connectRabbitMQ()
	if err != nil {
		logger.PrintFatal(err, nil, "general")
	}
	defer rabbitConn.Close()
	logger.PrintInfo("Соединение с RabbitMQ установлено", nil, "general")

	// monitoring
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config:      cfg,
		logger:      logger,
		models:      data.NewModels(db, redisClient),
		redis:       redisClient,
		rabbitMQ:    rabbitConn,
		mongoClient: mongoClient,
	}

	go app.consumeLogsFromRabbitMQ()

	if err = app.serve(); err != nil {
		logger.PrintFatal(err, nil, "general")
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}

// Подключение к Redis
func redisConnect() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "redis:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return redisClient, nil
}

// Подключение к MongoDB
func connectMongoDB() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoURI := getEnv("MONGO_URI", "mongodb://mongo:27017/qrent_logs")
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetAuth(options.Credential{
		Username:      getEnv("MONGO_USERNAME", "olzzhas"),
		Password:      getEnv("MONGO_PASSWORD", "Olzhas040404"),
		AuthMechanism: getEnv("MONGO_AUTH_MECHANISM", "SCRAM-SHA-256"),
	})
	clientOptions.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}

// Подключение к RabbitMQ
func connectRabbitMQ() (*amqp.Connection, error) {
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	return amqp.Dial(rabbitURL)
}

func (app *application) consumeLogsFromRabbitMQ() {
	ch, err := app.rabbitMQ.Channel()
	if err != nil {
		app.logger.PrintFatal(err, nil, "general")
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"logs", // Имя очереди
		"",     // Consumer tag
		true,   // Auto-ack
		false,  // Exclusive
		false,  // No-local
		false,  // No-wait
		nil,    // Arguments
	)
	if err != nil {
		app.logger.PrintFatal(err, nil, "general")
	}

	for msg := range msgs {
		var logMessage struct {
			Level      string         `json:"level"`
			Time       string         `json:"time"`
			Message    string         `json:"message"`
			Properties map[string]any `json:"properties,omitempty"`
			Collection string         `json:"collection"`
		}

		err := json.Unmarshal(msg.Body, &logMessage)
		if err != nil {
			app.logger.PrintError(err, nil, "general")
			continue
		}

		// Динамически выбираем коллекцию
		mongoCollection := app.mongoClient.Database("qrent_logs").Collection(logMessage.Collection)

		_, err = mongoCollection.InsertOne(context.TODO(), logMessage)
		if err != nil {
			app.logger.PrintError(err, nil, "general")
		}
	}
}
