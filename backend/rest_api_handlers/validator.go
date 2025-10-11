package rest_api_handlers

import "github.com/go-playground/validator/v10"

// validate is the shared validator instance used across all handlers
var validate *validator.Validate

func init() {
	validate = validator.New()
}
