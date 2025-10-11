package rest_api_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"family-calendar-backend/auth"
	"family-calendar-backend/db/services"
	"family-calendar-backend/rest_api_handlers/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// CreateCalendarMux creates a new calendar mux for the authenticated user
func CreateCalendarMux(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req CreateCalendarMuxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMsg := "Validation failed"
		if len(validationErrors) > 0 {
			errorMsg = utils.GetValidationErrorMsg(validationErrors[0])
		}
		utils.RespondError(w, http.StatusBadRequest, errorMsg, nil)
		return
	}

	calendarMux, err := services.CreateCalendarMux(userID, req.Name, req.Description)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create calendar mux", nil)
		return
	}

	// Build response
	response := CalendarMuxAPIResponse{
		ID:          calendarMux.ID,
		CreatedByID: calendarMux.CreatedByID,
		Name:        calendarMux.Name,
		Description: calendarMux.Description,
		CreatedAt:   calendarMux.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   calendarMux.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	utils.RespondJSON(w, http.StatusCreated, response)
}

// ListCalendarMuxes returns all calendar muxes owned by the authenticated user
func ListCalendarMuxes(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	calendarMuxes, err := services.GetCalendarMuxesByUser(userID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to retrieve calendar muxes", nil)
		return
	}

	// Build response
	calendarMuxResponses := make([]CalendarMuxAPIResponse, 0)
	for _, cm := range calendarMuxes {
		calendarMuxResponses = append(calendarMuxResponses, CalendarMuxAPIResponse{
			ID:          cm.ID,
			CreatedByID: cm.CreatedByID,
			Name:        cm.Name,
			Description: cm.Description,
			CreatedAt:   cm.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   cm.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	response := CalendarMuxListAPIResponse{
		CalendarMuxes: calendarMuxResponses,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// DeleteCalendarMux deletes a calendar mux owned by the authenticated user
func DeleteCalendarMux(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid calendar mux ID", nil)
		return
	}

	err = services.DeleteCalendarMux(uint(id), userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Calendar mux not found or access denied", nil)
		return
	}

	// Build response
	response := DeleteCalendarMuxAPIResponse{
		Message: "Calendar mux deleted successfully",
	}

	utils.RespondJSON(w, http.StatusOK, response)
}
