package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestRespondJSON(t *testing.T) {
	tests := []struct {
		name   string
		status int
		data   interface{}
		want   string
	}{
		{
			name:   "Simple map",
			status: http.StatusOK,
			data:   map[string]string{"message": "success"},
			want:   `{"message":"success"}`,
		},
		{
			name:   "Struct",
			status: http.StatusCreated,
			data:   struct{ ID int }{ID: 42},
			want:   `{"ID":42}`,
		},
		{
			name:   "Array",
			status: http.StatusOK,
			data:   []int{1, 2, 3},
			want:   `[1,2,3]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			RespondJSON(rr, tt.status, tt.data)

			assert.Equal(t, tt.status, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.want, rr.Body.String())
		})
	}
}

func TestRespondError(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		message string
		fields  map[string]string
		wantErr string
	}{
		{
			name:    "Simple error",
			status:  http.StatusBadRequest,
			message: "Invalid input",
			fields:  nil,
			wantErr: `{"error":"Invalid input"}`,
		},
		{
			name:    "Error with fields",
			status:  http.StatusBadRequest,
			message: "Validation failed",
			fields:  map[string]string{"email": "invalid format", "age": "too young"},
			wantErr: `{"error":"Validation failed","fields":{"email":"invalid format","age":"too young"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			RespondError(rr, tt.status, tt.message, tt.fields)

			assert.Equal(t, tt.status, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			assert.JSONEq(t, tt.wantErr, rr.Body.String())
		})
	}
}

func TestGetValidationErrorMsg(t *testing.T) {
	// Create a validator instance
	v := validator.New()

	// Create a struct with validation
	type TestStruct struct {
		Email string `validate:"required,email"`
		Age   int    `validate:"required,min=18,max=100"`
		Name  string `validate:"required"`
	}

	tests := []struct {
		name     string
		input    TestStruct
		wantMsgs map[string]string
	}{
		{
			name:  "Required field missing",
			input: TestStruct{Email: "", Age: 0, Name: ""},
			wantMsgs: map[string]string{
				"Email": "This field is required",
				"Age":   "This field is required",
				"Name":  "This field is required",
			},
		},
		{
			name:  "Email invalid format",
			input: TestStruct{Email: "not-an-email", Age: 25, Name: "John"},
			wantMsgs: map[string]string{
				"Email": "Invalid email format",
			},
		},
		{
			name:  "Age too small",
			input: TestStruct{Email: "john@example.com", Age: 10, Name: "John"},
			wantMsgs: map[string]string{
				"Age": "Value too small (min: 18)",
			},
		},
		{
			name:  "Age too large",
			input: TestStruct{Email: "john@example.com", Age: 150, Name: "John"},
			wantMsgs: map[string]string{
				"Age": "Value too large (max: 100)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.input)
			if err == nil {
				return
			}

			validationErrors := err.(validator.ValidationErrors)
			for _, fieldErr := range validationErrors {
				msg := GetValidationErrorMsg(fieldErr)
				expectedMsg, exists := tt.wantMsgs[fieldErr.Field()]
				if exists {
					assert.Equal(t, expectedMsg, msg)
				}
			}
		})
	}
}

func TestGetValidationErrorMsg_DefaultCase(t *testing.T) {
	// Create a validator instance
	v := validator.New()

	// Create a struct with a validation tag that isn't explicitly handled
	type TestStruct struct {
		URL string `validate:"url"`
	}

	input := TestStruct{URL: "not-a-url"}
	err := v.Struct(input)

	assert.Error(t, err)

	validationErrors := err.(validator.ValidationErrors)
	for _, fieldErr := range validationErrors {
		msg := GetValidationErrorMsg(fieldErr)
		assert.Equal(t, "Invalid value", msg)
	}
}
