package db

import (
	"family-calendar-backend/db/models"
	"flag"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// migrateFunc allows mocking AutoMigrate in tests
var migrateFunc = func(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.CalendarMux{})
}

// getSQLitePath returns the appropriate SQLite database path
// based on whether we're running tests or in production
var getSQLitePath = func() string {
	if flag.Lookup("test.v") != nil {
		// Running tests - use in-memory database
		return ":memory:"
	}
	// Production - use default file
	return "family_calendar.db"
}

func InitDB() error {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // Default to SQLite
	}

	var err error
	var dialector gorm.Dialector

	switch dbType {
	case "postgres":
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			// Fallback to individual postgres environment variables
			host := os.Getenv("DB_HOST")
			if host == "" {
				host = "localhost"
			}
			port := os.Getenv("DB_PORT")
			if port == "" {
				port = "5432"
			}
			user := os.Getenv("DB_USER")
			if user == "" {
				user = "postgres"
			}
			password := os.Getenv("DB_PASSWORD")
			if password == "" {
				password = "postgres"
			}
			dbname := os.Getenv("DB_NAME")
			if dbname == "" {
				dbname = "family_calendar"
			}
			sslmode := os.Getenv("DB_SSLMODE")
			if sslmode == "" {
				sslmode = "disable"
			}

			dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
				host, port, user, password, dbname, sslmode)
		}
		dialector = postgres.Open(dsn)

	case "sqlite":
		dialector = sqlite.Open(getSQLitePath())

	default:
		return fmt.Errorf("unsupported database type: %s (supported: sqlite, postgres)", dbType)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
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
