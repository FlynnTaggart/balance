package main

import (
	"balance/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")))
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
