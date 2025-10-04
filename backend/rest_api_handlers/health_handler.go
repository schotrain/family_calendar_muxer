package rest_api_handlers

import (
	"net/http"

	"family-calendar-backend/rest_api_handlers/utils"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
