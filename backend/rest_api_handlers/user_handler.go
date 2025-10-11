package rest_api_handlers

import (
	"net/http"

	"family-calendar-backend/auth"
	"family-calendar-backend/db"
	"family-calendar-backend/db/models"
	"family-calendar-backend/rest_api_handlers/utils"
)

func UserInfo(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user ID from context
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Query user from database
	var dbUser models.User
	if err := db.DB.First(&dbUser, userID).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, "User not found", nil)
		return
	}

	// Build response
	response := UserAPIResponse{
		ID:         int(dbUser.ID),
		GivenName:  dbUser.GivenName,
		FamilyName: dbUser.FamilyName,
		Email:      dbUser.Email,
	}

	// Validate response
	if err := validate.Struct(response); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Response validation failed", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, response)
}
