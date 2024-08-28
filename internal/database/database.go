package db

import (
	"log"
	"test-task/internal/modules/auth/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

type DBHandler struct {
	DB *gorm.DB
}

func InitDB(url string) DBHandler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalf("error occured while opening db conn: %s", err)
	}

	log.Println("database connected: ", url)

	err = db.AutoMigrate(
		&models.User{},
		&models.Token{},
	)
	if err != nil {
		log.Fatalf("error occurred while migration: %s", err)
	}

	log.Println("migrations done")

	dbConn = db
	return DBHandler{db}
}

func GetDBHandler() *DBHandler {
	if dbConn == nil {
		log.Fatal("database connection is not initialized")
	}

	return &DBHandler{DB: dbConn}
}
