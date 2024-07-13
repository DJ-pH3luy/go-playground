package db

import (
	"log"

	"github.com/dj-ph3luy/go-playground/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct{
	ConnectionString string
}

func Connect(conf Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(conf.ConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func Migrate(database *gorm.DB) {
	database.AutoMigrate(&entities.User{})
}