package db

import (
	"family-calendar-backend/db/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// migrateFunc allows mocking AutoMigrate in tests
var migrateFunc = func(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{})
}

func InitDB(dbPath ...string) error {
	// Use provided path or default to "family_calendar.db"
	path := "family_calendar.db"
	if len(dbPath) > 0 && dbPath[0] != "" {
		path = dbPath[0]
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return err
	}

	// Run migrations
	err = migrateFunc(DB)
	if err != nil {
		return err
	}

	return nil
}
