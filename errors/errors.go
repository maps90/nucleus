package errors

import (
	"net/http"
	"sort"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
)

type validationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// Errors struct
type Errors struct {
	ctx echo.Context
}

// InternalServerError creates a new API error representing an internal server error (HTTP 500)
func (e *Errors) InternalServerError(err error) *APIError {
	apiError := NewAPIError(e.ctx, "Internal Server Error", Params{"error": err})
	apiError.Status = http.StatusInternalServerError
	return apiError
}

// NotFound creates a new API error representing a resource-not-found error (HTTP 404)
func (e *Errors) NotFound(resource string) *APIError {
	apiError := NewAPIError(e.ctx, "Not Found", Params{"resource": resource})
	apiError.Status = http.StatusNotFound
	return apiError
}

// Unauthorized creates a new API error representing an authentication failure (HTTP 401)
func (e *Errors) Unauthorized(err string) *APIError {
	apiError := NewAPIError(e.ctx, "Unauthorized", Params{"error": err})
	apiError.Status = http.StatusUnauthorized
	return apiError
}

// InvalidData converts a data validation error into an API error (HTTP 400)
func (e *Errors) InvalidData(errs validation.Errors) *APIError {
	result := make([]validationError, 0)
	fields := make([]string, 0)
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		err := errs[field]
		result = append(result, validationError{
			Field: field,
			Error: err.Error(),
		})
	}

	err := NewAPIError(e.ctx, "Invalid Data", nil)
	err.Status = http.StatusBadRequest
	err.Details = result

	return err
}
