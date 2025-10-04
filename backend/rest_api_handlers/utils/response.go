package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string, fields map[string]string) {
	RespondJSON(w, status, ErrorResponse{
		Error:  message,
		Fields: fields,
	})
}

func GetValidationErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value too small (min: " + err.Param() + ")"
	case "max":
		return "Value too large (max: " + err.Param() + ")"
	default:
		return "Invalid value"
	}
}
