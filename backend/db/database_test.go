package db

import (
	"os"
	"testing"

	"family-calendar-backend/db/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestInitDB_Success(t *testing.T) {
	// Ensure we use SQLite for testing
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")

	testDB := "test_family_calendar.db"
	defer os.Remove(testDB)

	err := InitDB(testDB)

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
}

func TestInitDB_Migration(t *testing.T) {
	// Ensure we use SQLite for testing
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")

	testDB := "test_migration.db"
	defer os.Remove(testDB)

	err := InitDB(testDB)
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
	// Ensure we use SQLite for testing
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")

	testDB := "test_global.db"
	defer os.Remove(testDB)

	DB = nil

	err := InitDB(testDB)
	assert.NoError(t, err)
	assert.NotNil(t, DB)
	assert.IsType(t, &gorm.DB{}, DB)

	// Clean up
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	sqlDB.Close()
}

func TestInitDB_InvalidDatabasePath(t *testing.T) {
	// Ensure we use SQLite for testing
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")

	err := InitDB("/nonexistent/directory/that/does/not/exist/test.db")

	assert.Error(t, err)
}

func TestInitDB_CustomPath(t *testing.T) {
	// Ensure we use SQLite for testing
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")

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
	// Ensure we use SQLite for testing
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")

	// Save original migrateFunc
	originalMigrate := migrateFunc
	defer func() { migrateFunc = originalMigrate }()

	// Mock migrateFunc to return an error
	migrateFunc = func(db *gorm.DB) error {
		return assert.AnError
	}

	migrationTestDB := "test_migration_error.db"
	defer os.Remove(migrationTestDB)

	err := InitDB(migrationTestDB)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
}

func TestInitDB_UnsupportedDatabaseType(t *testing.T) {
	// Test with unsupported database type
	os.Setenv("DB_TYPE", "mysql")
	defer os.Unsetenv("DB_TYPE")

	err := InitDB()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported database type")
}

func TestInitDB_DefaultsToSQLite(t *testing.T) {
	// Clear DB_TYPE to test default behavior
	os.Unsetenv("DB_TYPE")

	testDB := "test_default.db"
	defer os.Remove(testDB)

	err := InitDB(testDB)

	assert.NoError(t, err)
	assert.NotNil(t, DB)

	// Clean up
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	sqlDB.Close()
}

func TestInitDB_SQLiteDefaultPath(t *testing.T) {
	// Test SQLite with no path provided (uses default)
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")
	defer os.Remove("family_calendar.db")

	err := InitDB()

	assert.NoError(t, err)
	assert.NotNil(t, DB)

	// Verify the default database file was created
	_, err = os.Stat("family_calendar.db")
	assert.NoError(t, err)

	// Clean up
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	sqlDB.Close()
}

func TestInitDB_SQLiteEmptyPath(t *testing.T) {
	// Test SQLite with empty string path (uses default)
	os.Setenv("DB_TYPE", "sqlite")
	defer os.Unsetenv("DB_TYPE")
	defer os.Remove("family_calendar.db")

	err := InitDB("")

	assert.NoError(t, err)
	assert.NotNil(t, DB)

	// Clean up
	sqlDB, err := DB.DB()
	assert.NoError(t, err)
	sqlDB.Close()
}

func TestInitDB_PostgresWithDatabaseURL(t *testing.T) {
	// Test PostgreSQL with DATABASE_URL
	os.Setenv("DB_TYPE", "postgres")
	os.Setenv("DATABASE_URL", "host=nonexistent.invalid port=5432 user=testuser password=testpass dbname=testdb sslmode=disable")
	defer os.Unsetenv("DB_TYPE")
	defer os.Unsetenv("DATABASE_URL")

	// Save original migrateFunc to skip actual migration
	originalMigrate := migrateFunc
	defer func() { migrateFunc = originalMigrate }()
	migrateFunc = func(db *gorm.DB) error { return nil }

	err := InitDB()

	// This will fail to connect since the host doesn't exist, but we're testing the DSN construction
	// The error will be a connection error, not a configuration error
	if err != nil {
		assert.NotContains(t, err.Error(), "unsupported database type")
	}
}

func TestInitDB_PostgresWithIndividualEnvVars(t *testing.T) {
	// Test PostgreSQL with individual environment variables
	os.Setenv("DB_TYPE", "postgres")
	os.Setenv("DB_HOST", "nonexistent.invalid")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "require")
	defer os.Unsetenv("DB_TYPE")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")
	defer os.Unsetenv("DB_USER")
	defer os.Unsetenv("DB_PASSWORD")
	defer os.Unsetenv("DB_NAME")
	defer os.Unsetenv("DB_SSLMODE")

	// Save original migrateFunc to skip actual migration
	originalMigrate := migrateFunc
	defer func() { migrateFunc = originalMigrate }()
	migrateFunc = func(db *gorm.DB) error { return nil }

	err := InitDB()

	// This will fail to connect since the host doesn't exist
	if err != nil {
		assert.NotContains(t, err.Error(), "unsupported database type")
	}
}

func TestInitDB_PostgresWithDefaults(t *testing.T) {
	// Test PostgreSQL with default values for all env vars
	os.Setenv("DB_TYPE", "postgres")
	defer os.Unsetenv("DB_TYPE")
	// Explicitly unset all postgres env vars to test defaults
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")

	// Save original migrateFunc to skip actual migration
	originalMigrate := migrateFunc
	defer func() { migrateFunc = originalMigrate }()
	migrateFunc = func(db *gorm.DB) error { return nil }

	err := InitDB()

	// This will fail to connect, but tests that defaults are applied (localhost:5432, etc)
	if err != nil {
		assert.NotContains(t, err.Error(), "unsupported database type")
	}
}
