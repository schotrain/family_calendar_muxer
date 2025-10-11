package services

import (
	"errors"
	"regexp"
	"testing"

	"family-calendar-backend/db"
	"family-calendar-backend/db/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	var err error
	db.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.DB.AutoMigrate(&models.User{}, &models.CalendarMux{})
	assert.NoError(t, err)
}

func TestCreateCalendarMux(t *testing.T) {
	setupTestDB(t)

	// Create a test user first
	user := models.User{
		GivenName:      "Test",
		FamilyName:     "User",
		Email:          "test@example.com",
		AuthProvider:   "google",
		AuthProviderID: "test-123",
	}
	db.DB.Create(&user)

	// Test creating a calendar mux
	calendarMux, err := CreateCalendarMux(user.ID, "Test Calendar", "Test Description")

	assert.NoError(t, err)
	assert.NotNil(t, calendarMux)
	assert.Equal(t, "Test Calendar", calendarMux.Name)
	assert.Equal(t, "Test Description", calendarMux.Description)
	assert.Equal(t, user.ID, calendarMux.CreatedByID)
	assert.NotZero(t, calendarMux.ID)
}

func TestGetCalendarMuxesByUser(t *testing.T) {
	setupTestDB(t)

	// Create test users
	user1 := models.User{
		GivenName:      "User",
		FamilyName:     "One",
		Email:          "user1@example.com",
		AuthProvider:   "google",
		AuthProviderID: "user1-123",
	}
	db.DB.Create(&user1)

	user2 := models.User{
		GivenName:      "User",
		FamilyName:     "Two",
		Email:          "user2@example.com",
		AuthProvider:   "google",
		AuthProviderID: "user2-123",
	}
	db.DB.Create(&user2)

	// Create calendar muxes for user1
	db.DB.Create(&models.CalendarMux{
		CreatedByID: user1.ID,
		Name:        "User 1 Calendar 1",
		Description: "First calendar",
	})
	db.DB.Create(&models.CalendarMux{
		CreatedByID: user1.ID,
		Name:        "User 1 Calendar 2",
		Description: "Second calendar",
	})

	// Create calendar mux for user2
	db.DB.Create(&models.CalendarMux{
		CreatedByID: user2.ID,
		Name:        "User 2 Calendar",
		Description: "Other user calendar",
	})

	// Test getting calendar muxes for user1
	calendarMuxes, err := GetCalendarMuxesByUser(user1.ID)

	assert.NoError(t, err)
	assert.Len(t, calendarMuxes, 2)
	assert.Equal(t, "User 1 Calendar 1", calendarMuxes[0].Name)
	assert.Equal(t, "User 1 Calendar 2", calendarMuxes[1].Name)

	// Test getting calendar muxes for user2
	calendarMuxes, err = GetCalendarMuxesByUser(user2.ID)

	assert.NoError(t, err)
	assert.Len(t, calendarMuxes, 1)
	assert.Equal(t, "User 2 Calendar", calendarMuxes[0].Name)
}

func TestGetCalendarMuxesByUser_EmptyList(t *testing.T) {
	setupTestDB(t)

	// Create a user with no calendar muxes
	user := models.User{
		GivenName:      "Test",
		FamilyName:     "User",
		Email:          "test@example.com",
		AuthProvider:   "google",
		AuthProviderID: "test-123",
	}
	db.DB.Create(&user)

	// Test getting calendar muxes
	calendarMuxes, err := GetCalendarMuxesByUser(user.ID)

	assert.NoError(t, err)
	assert.Len(t, calendarMuxes, 0)
}

func TestDeleteCalendarMux_Success(t *testing.T) {
	setupTestDB(t)

	// Create a test user
	user := models.User{
		GivenName:      "Test",
		FamilyName:     "User",
		Email:          "test@example.com",
		AuthProvider:   "google",
		AuthProviderID: "test-123",
	}
	db.DB.Create(&user)

	// Create a calendar mux
	calendarMux := models.CalendarMux{
		CreatedByID: user.ID,
		Name:        "Test Calendar",
		Description: "Test Description",
	}
	db.DB.Create(&calendarMux)

	// Test deleting the calendar mux
	err := DeleteCalendarMux(calendarMux.ID, user.ID)

	assert.NoError(t, err)

	// Verify it was deleted
	var found models.CalendarMux
	result := db.DB.First(&found, calendarMux.ID)
	assert.Error(t, result.Error)
}

func TestDeleteCalendarMux_WrongUser(t *testing.T) {
	setupTestDB(t)

	// Create two users
	user1 := models.User{
		GivenName:      "User",
		FamilyName:     "One",
		Email:          "user1@example.com",
		AuthProvider:   "google",
		AuthProviderID: "user1-123",
	}
	db.DB.Create(&user1)

	user2 := models.User{
		GivenName:      "User",
		FamilyName:     "Two",
		Email:          "user2@example.com",
		AuthProvider:   "google",
		AuthProviderID: "user2-123",
	}
	db.DB.Create(&user2)

	// Create a calendar mux for user1
	calendarMux := models.CalendarMux{
		CreatedByID: user1.ID,
		Name:        "User 1 Calendar",
		Description: "Test Description",
	}
	db.DB.Create(&calendarMux)

	// Try to delete with user2 (should fail)
	err := DeleteCalendarMux(calendarMux.ID, user2.ID)

	assert.Error(t, err)

	// Verify it was NOT deleted
	var found models.CalendarMux
	result := db.DB.First(&found, calendarMux.ID)
	assert.NoError(t, result.Error)
}

func TestDeleteCalendarMux_NotFound(t *testing.T) {
	setupTestDB(t)

	// Create a test user
	user := models.User{
		GivenName:      "Test",
		FamilyName:     "User",
		Email:          "test@example.com",
		AuthProvider:   "google",
		AuthProviderID: "test-123",
	}
	db.DB.Create(&user)

	// Try to delete a non-existent calendar mux
	err := DeleteCalendarMux(9999, user.ID)

	assert.Error(t, err)
}

// Error scenario tests using sqlmock
func TestCreateCalendarMux_DatabaseError(t *testing.T) {
	// Create mock DB
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Set up the global DB
	originalDB := db.DB
	db.DB = gormDB
	defer func() { db.DB = originalDB }()

	// Expect the INSERT to fail
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "calendar_muxes"`)).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	// Test creating a calendar mux
	_, err = CreateCalendarMux(1, "Test", "Description")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestGetCalendarMuxesByUser_DatabaseError(t *testing.T) {
	// Create mock DB
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Set up the global DB
	originalDB := db.DB
	db.DB = gormDB
	defer func() { db.DB = originalDB }()

	// Expect the SELECT to fail
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "calendar_muxes"`)).
		WillReturnError(errors.New("database error"))

	// Test getting calendar muxes
	_, err = GetCalendarMuxesByUser(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestDeleteCalendarMux_DatabaseError(t *testing.T) {
	// Create mock DB
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Set up the global DB
	originalDB := db.DB
	db.DB = gormDB
	defer func() { db.DB = originalDB }()

	// Expect the DELETE to fail
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "calendar_muxes"`)).
		WillReturnError(errors.New("database error"))
	mock.ExpectRollback()

	// Test deleting a calendar mux
	err = DeleteCalendarMux(1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}
