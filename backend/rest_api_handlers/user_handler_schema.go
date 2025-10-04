package rest_api_handlers

type CreateUserAPIRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=1,max=150"`
}

type UserAPIResponse struct {
	ID    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=1,max=150"`
}
