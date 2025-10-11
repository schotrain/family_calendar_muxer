package rest_api_handlers

type CreateCalendarMuxRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=200"`
	Description string `json:"description" validate:"max=1000"`
}

type CalendarMuxAPIResponse struct {
	ID          uint   `json:"id" validate:"required"`
	CreatedByID uint   `json:"created_by_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=1,max=200"`
	Description string `json:"description" validate:"max=1000"`
	CreatedAt   string `json:"created_at" validate:"required"`
	UpdatedAt   string `json:"updated_at" validate:"required"`
}

type CalendarMuxListAPIResponse struct {
	CalendarMuxes []CalendarMuxAPIResponse `json:"calendar_muxes" validate:"dive"`
}

type DeleteCalendarMuxAPIResponse struct {
	Message string `json:"message" validate:"required"`
}
