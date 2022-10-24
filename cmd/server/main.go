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
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := fiber.New()

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
	)

	//db, err := gorm.Open(postgres.Open(dsn))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//sqlDB, err := db.DB()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//sqlDB.SetMaxIdleConns(10)
	//sqlDB.SetMaxOpenConns(100)
	//sqlDB.SetConnMaxLifetime(time.Hour)
	//
	//postgresDB := databases.NewPostgresDB(db)
	//

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}
	config.MaxConns = 50
	config.MaxConnLifetime = time.Minute * 10
	config.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.Connect(context.Background(), dsn)

	pgxDB := databases.NewPgxDB(pool)

	handler := handlers.NewHandler(pgxDB)

	routes.InitializeRoutes(app, handler)

	if err = app.Listen(os.Getenv("SERVER_URL")); err != nil {
		log.Printf("Server is not running! Reason: %v", err)
	}
}
