package rest_api_handlers

type UserAPIResponse struct {
	ID         int    `json:"id" validate:"required"`
	GivenName  string `json:"given_name" validate:"required,min=2,max=100"`
	FamilyName string `json:"family_name" validate:"required,min=2,max=100"`
	Email      string `json:"email" validate:"required,email"`
}
