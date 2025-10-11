package rest_api_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"family-calendar-backend/auth"
	"family-calendar-backend/db"
	"family-calendar-backend/db/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupCalendarMuxTestDB(t *testing.T) *models.User {
	var err error
	db.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.DB.AutoMigrate(&models.User{}, &models.CalendarMux{})
	assert.NoError(t, err)

	// Create a test user
	user := &models.User{
		GivenName:      "Test",
		FamilyName:     "User",
		Email:          "test@example.com",
		AuthProvider:   "google",
		AuthProviderID: "test-123",
	}
	db.DB.Create(user)

	return user
}

func TestCreateCalendarMux_Success(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	reqBody := CreateCalendarMuxRequest{
		Name:        "Test Calendar",
		Description: "Test Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/calendar-mux", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	CreateCalendarMux(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response CalendarMuxAPIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Test Calendar", response.Name)
	assert.Equal(t, "Test Description", response.Description)
	assert.Equal(t, user.ID, response.CreatedByID)
}

func TestCreateCalendarMux_NoAuth(t *testing.T) {
	setupCalendarMuxTestDB(t)

	reqBody := CreateCalendarMuxRequest{
		Name:        "Test Calendar",
		Description: "Test Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/calendar-mux", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	CreateCalendarMux(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCreateCalendarMux_InvalidJSON(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	req := httptest.NewRequest(http.MethodPost, "/api/calendar-mux", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	CreateCalendarMux(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateCalendarMux_ValidationError(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	// Name is required
	reqBody := CreateCalendarMuxRequest{
		Name:        "",
		Description: "Test Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/calendar-mux", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	CreateCalendarMux(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestListCalendarMuxes_Success(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	// Create test calendar muxes
	db.DB.Create(&models.CalendarMux{
		CreatedByID: user.ID,
		Name:        "Calendar 1",
		Description: "Description 1",
	})
	db.DB.Create(&models.CalendarMux{
		CreatedByID: user.ID,
		Name:        "Calendar 2",
		Description: "Description 2",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/calendar-mux", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	ListCalendarMuxes(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response CalendarMuxListAPIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.CalendarMuxes, 2)
	assert.Equal(t, "Calendar 1", response.CalendarMuxes[0].Name)
	assert.Equal(t, "Calendar 2", response.CalendarMuxes[1].Name)
}

func TestListCalendarMuxes_EmptyList(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/calendar-mux", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	ListCalendarMuxes(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response CalendarMuxListAPIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.CalendarMuxes, 0)
}

func TestListCalendarMuxes_NoAuth(t *testing.T) {
	setupCalendarMuxTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/calendar-mux", nil)

	rr := httptest.NewRecorder()

	ListCalendarMuxes(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeleteCalendarMux_Success(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	// Create a calendar mux
	calendarMux := &models.CalendarMux{
		CreatedByID: user.ID,
		Name:        "Test Calendar",
		Description: "Test Description",
	}
	db.DB.Create(calendarMux)

	req := httptest.NewRequest(http.MethodDelete, "/api/calendar-mux/1", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)

	// Add chi URL param
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	DeleteCalendarMux(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response DeleteCalendarMuxAPIResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Calendar mux deleted successfully", response.Message)

	// Verify it was deleted
	var found models.CalendarMux
	result := db.DB.First(&found, calendarMux.ID)
	assert.Error(t, result.Error)
}

func TestDeleteCalendarMux_NoAuth(t *testing.T) {
	setupCalendarMuxTestDB(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/calendar-mux/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	DeleteCalendarMux(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestDeleteCalendarMux_InvalidID(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/calendar-mux/invalid", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	DeleteCalendarMux(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteCalendarMux_NotFound(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	req := httptest.NewRequest(http.MethodDelete, "/api/calendar-mux/9999", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "9999")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	DeleteCalendarMux(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteCalendarMux_WrongUser(t *testing.T) {
	user := setupCalendarMuxTestDB(t)

	// Create another user
	otherUser := &models.User{
		GivenName:      "Other",
		FamilyName:     "User",
		Email:          "other@example.com",
		AuthProvider:   "google",
		AuthProviderID: "other-123",
	}
	db.DB.Create(otherUser)

	// Create a calendar mux for the other user
	calendarMux := &models.CalendarMux{
		CreatedByID: otherUser.ID,
		Name:        "Other's Calendar",
		Description: "Test Description",
	}
	db.DB.Create(calendarMux)

	// Try to delete as first user
	req := httptest.NewRequest(http.MethodDelete, "/api/calendar-mux/1", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, user.ID)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	DeleteCalendarMux(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Verify it was NOT deleted
	var found models.CalendarMux
	result := db.DB.First(&found, calendarMux.ID)
	assert.NoError(t, result.Error)
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

	reqBody := CreateCalendarMuxRequest{
		Name:        "Test Calendar",
		Description: "Test Description",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/calendar-mux", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, uint(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	CreateCalendarMux(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestListCalendarMuxes_DatabaseError(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/api/calendar-mux", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDContextKey, uint(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	ListCalendarMuxes(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
