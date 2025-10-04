package rest_api_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"family-calendar-backend/database"
	"family-calendar-backend/db_models"
	"family-calendar-backend/rest_api_handlers/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserAPIRequest

	// Decode request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid JSON", nil)
		return
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = utils.GetValidationErrorMsg(err)
		}
		utils.RespondError(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	// Create user in database
	dbUser := db_models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	if err := database.DB.Create(&dbUser).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create user", nil)
		return
	}

	// Build response
	response := UserAPIResponse{
		ID:    int(dbUser.ID),
		Name:  dbUser.Name,
		Email: dbUser.Email,
		Age:   dbUser.Age,
	}

	// Validate response
	if err := validate.Struct(response); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Response validation failed", nil)
		return
	}

	utils.RespondJSON(w, http.StatusCreated, response)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL parameter
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	// Query user from database
	var dbUser db_models.User
	if err := database.DB.First(&dbUser, id).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, "User not found", nil)
		return
	}

	// Build response
	response := UserAPIResponse{
		ID:    int(dbUser.ID),
		Name:  dbUser.Name,
		Email: dbUser.Email,
		Age:   dbUser.Age,
	}

	// Validate response
	if err := validate.Struct(response); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Response validation failed", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, response)
}
