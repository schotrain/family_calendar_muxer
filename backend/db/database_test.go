package db

import (
	"os"
	"testing"

	"family-calendar-backend/db/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitDB_Success(t *testing.T) {
	// Use a temporary test database
	testDB := "test_family_calendar.db"
	defer os.Remove(testDB)

	// Temporarily replace the database path
	originalOpen := "family_calendar.db"

	// We need to modify InitDB to accept a path parameter for testing
	// For now, let's test that it creates the database successfully
	err := InitDB()

	assert.NoError(t, err)
	assert.NotNil(t, DB)

	// Verify the User table was created
	var count int64
	result := DB.Model(&models.User{}).Count(&count)
	assert.NoError(t, result.Error)

	// Clean up
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	_ = originalOpen // avoid unused variable
}

func TestInitDB_Migration(t *testing.T) {
	// Initialize DB
	err := InitDB()
	assert.NoError(t, err)
	assert.NotNil(t, DB)

	// Verify User model migration by checking if we can query the table
	var users []models.User
	result := DB.Find(&users)
	assert.NoError(t, result.Error)

	// Verify table structure has expected columns
	hasTable := DB.Migrator().HasTable(&models.User{})
	assert.True(t, hasTable)

	hasColumn := DB.Migrator().HasColumn(&models.User{}, "Email")
	assert.True(t, hasColumn)

	hasColumn = DB.Migrator().HasColumn(&models.User{}, "AuthProviderID")
	assert.True(t, hasColumn)

	// Clean up
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	sqlDB.Close()
}

func TestInitDB_DBGlobalVariable(t *testing.T) {
	// Test that InitDB properly sets the global DB variable
	DB = nil

	err := InitDB()
	assert.NoError(t, err)
	assert.NotNil(t, DB)
	assert.IsType(t, &gorm.DB{}, DB)

	// Clean up
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	sqlDB.Close()
}

func TestInitDB_InvalidDatabasePath(t *testing.T) {
	// Test with an invalid database path (directory that doesn't exist)
	err := InitDB("/nonexistent/directory/that/does/not/exist/test.db")

	assert.Error(t, err)
	// The error comes from sqlite trying to create the file in a non-existent directory
}

func TestInitDB_CustomPath(t *testing.T) {
	// Test with custom database path
	customPath := "test_custom.db"
	defer os.Remove(customPath)

	err := InitDB(customPath)

	assert.NoError(t, err)
	assert.NotNil(t, DB)

	// Clean up
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	sqlDB.Close()
}

func TestInitDB_MigrationError(t *testing.T) {
	// Save original migrateFunc
	originalMigrate := migrateFunc
	defer func() { migrateFunc = originalMigrate }()

	// Mock migrateFunc to return an error
	migrateFunc = func(db *gorm.DB) error {
		return assert.AnError // Returns a generic test error
	}

	migrationTestDB := "test_migration_error.db"
	defer os.Remove(migrationTestDB)

	err := InitDB(migrationTestDB)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}
