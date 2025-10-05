package rest_api_handlers

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"family-calendar-backend/auth"
	"family-calendar-backend/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	db.DB = gormDB

	cleanup := func() {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}

	return mock, cleanup
}

func TestUserInfo_Success(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	now := time.Now()
	// Mock database query - GORM adds deleted_at check, ORDER BY, and LIMIT
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "given_name", "family_name", "email", "auth_provider", "auth_provider_id"}).
		AddRow(123, now, now, nil, "John", "Doe", "john@example.com", "google", "google-123")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(123, 1).
		WillReturnRows(rows)

	// Create request with user ID in context
	req := httptest.NewRequest("GET", "/api/userinfo", nil)
	ctx := req.Context()
	ctx = auth.SetUserIDInContext(ctx, 123)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	UserInfo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"id":123`)
	assert.Contains(t, rr.Body.String(), `"given_name":"John"`)
	assert.Contains(t, rr.Body.String(), `"family_name":"Doe"`)
	assert.Contains(t, rr.Body.String(), `"email":"john@example.com"`)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserInfo_NoUserInContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/userinfo", nil)
	rr := httptest.NewRecorder()

	UserInfo(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "User not authenticated")
}

func TestUserInfo_UserNotFound(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	// Mock database query returning no results - GORM adds deleted_at check, ORDER BY, and LIMIT
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Create request with user ID in context
	req := httptest.NewRequest("GET", "/api/userinfo", nil)
	ctx := req.Context()
	ctx = auth.SetUserIDInContext(ctx, 999)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	UserInfo(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "User not found")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserInfo_ValidationFailure(t *testing.T) {
	mock, cleanup := setupTestDB(t)
	defer cleanup()

	now := time.Now()
	// Mock database query with invalid data (empty required fields) to trigger validation failure
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "given_name", "family_name", "email", "auth_provider", "auth_provider_id"}).
		AddRow(456, now, now, nil, "", "", "", "google", "google-456") // Empty required fields

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(456, 1).
		WillReturnRows(rows)

	// Create request with user ID in context
	req := httptest.NewRequest("GET", "/api/userinfo", nil)
	ctx := req.Context()
	ctx = auth.SetUserIDInContext(ctx, 456)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	UserInfo(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Response validation failed")

	assert.NoError(t, mock.ExpectationsWereMet())
}
