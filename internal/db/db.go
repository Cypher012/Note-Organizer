package db

import (
	"log"

	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("storage.db"))
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	if err := db.AutoMigrate(&models.User{}, &models.Folder{}, &models.Note{}); err != nil {
		log.Fatal("Failed to migrate database")
	}

	return db
}
