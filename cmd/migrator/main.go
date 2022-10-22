package main

import (
	"balance/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	dsn := fmt.Sprintf("host=localhost user=%v password=%v dbname=%v port=%v sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&models.Service{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&models.Report{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&models.Reserve{})
	if err != nil {
		log.Fatalln(err)
	}

}
