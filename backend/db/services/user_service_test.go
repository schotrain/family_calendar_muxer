package services

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"family-calendar-backend/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	// Set the global DB
	db.DB = gormDB

	cleanup := func() {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}

	return mock, cleanup
}

func TestFindOrCreateUser_UserExists(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()
	// Mock finding existing user - GORM adds deleted_at check, ORDER BY, and LIMIT
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "given_name", "family_name", "email", "auth_provider", "auth_provider_id"}).
		AddRow(1, now, now, nil, "John", "Doe", "john@example.com", "google", "google-123")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (auth_provider = $1 AND auth_provider_id = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
		WithArgs("google", "google-123", 1).
		WillReturnRows(rows)

	// Mock update
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	user, err := FindOrCreateUser("google", "google-123", "Jane", "Smith", "jane@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "Jane", user.GivenName)
	assert.Equal(t, "Smith", user.FamilyName)
	assert.Equal(t, "jane@example.com", user.Email)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindOrCreateUser_CreateNewUser(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	// Mock user not found - GORM adds deleted_at check, ORDER BY, and LIMIT
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (auth_provider = $1 AND auth_provider_id = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
		WithArgs("google", "google-456", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Mock create
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			"Alice",
			"Johnson",
			"alice@example.com",
			"google",
			"google-456",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	mock.ExpectCommit()

	user, err := FindOrCreateUser("google", "google-456", "Alice", "Johnson", "alice@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Alice", user.GivenName)
	assert.Equal(t, "Johnson", user.FamilyName)
	assert.Equal(t, "alice@example.com", user.Email)
	assert.Equal(t, "google", user.AuthProvider)
	assert.Equal(t, "google-456", user.AuthProviderID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindOrCreateUser_CreateError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()

	// Mock user not found - GORM adds deleted_at check, ORDER BY, and LIMIT
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE (auth_provider = $1 AND auth_provider_id = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`)).
		WithArgs("google", "google-789", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Mock create error
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()

	user, err := FindOrCreateUser("google", "google-789", "Bob", "Brown", "bob@example.com")

	assert.Error(t, err)
	assert.Nil(t, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}
