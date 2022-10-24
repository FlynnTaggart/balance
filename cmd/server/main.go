package main

import (
	"balance/internal/databases"
	"balance/internal/handlers"
	"balance/internal/routes"
	"context"
	"fmt"

	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initializeLogger() *zap.Logger {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile("./logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.InfoLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func main() {
	app := fiber.New()

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"))

	logger := initializeLogger()
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("zap: %v", err)
		}
	}()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}
	config.ConnConfig.Logger = zapadapter.NewLogger(logger)
	config.ConnConfig.LogLevel = pgx.LogLevelInfo
	config.MaxConns = 50
	config.MaxConnLifetime = time.Minute * 10
	config.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.ConnectConfig(context.Background(), config)

	pgxDB := databases.NewPgxDB(pool)

	handler := handlers.NewHandler(pgxDB)

	routes.InitializeRoutes(app, handler)

	if err = app.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Server is not running! Reason: %v", err)
	}
}
