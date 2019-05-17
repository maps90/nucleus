package errors

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/maps90/nucleus/config"
)

type (
	// Params is used to replace placeholders in an error template with the corresponding values.
	Params map[string]interface{}

	errorTemplate struct {
		Message          string `yaml:"message"`
		DeveloperMessage string `yaml:"developer_message"`
	}
)

var templates = map[string]errorTemplate{
	"Internal Server Error": errorTemplate{
		Message:          "We have encountered an internal server error.",
		DeveloperMessage: "Internal Server Error: {error}",
	},
	"Not Found": errorTemplate{
		Message: "{resource} was not found.",
	},

	"Unauthorized": errorTemplate{
		Message:          "Authentication failed.",
		DeveloperMessage: "Authentication failed: {error}",
	},

	"Invalid Data": errorTemplate{
		Message: "There is some problem with the data you submitted. See \"details\" for more information.",
	},

	"Bad Request": errorTemplate{
		Message:          "Missing required parameter.",
		DeveloperMessage: "Bad Request: {error}",
	},
}

// NewAPIError creates a new APIError with the given HTTP status code, error code, and parameters for replacing placeholders in the error template.
// The param can be nil, indicating there is no need for placeholder replacement.
func NewAPIError(c echo.Context, code string, params Params) *APIError {
	err := &APIError{
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
		ErrorCode: code,
		Message:   code,
	}

	if template, ok := templates[code]; ok {
		err.Message = template.getMessage(params)
		if config.Get("env") != "production" {
			err.DeveloperMessage = template.getDeveloperMessage(params)
		}
	}

	return err
}

// getMessage returns the error message by replacing placeholders in the error template with the actual parameters.
func (e errorTemplate) getMessage(params Params) string {
	return replacePlaceholders(e.Message, params)
}

// getDeveloperMessage returns the developer message by replacing placeholders in the error template with the actual parameters.
func (e errorTemplate) getDeveloperMessage(params Params) string {
	return replacePlaceholders(e.DeveloperMessage, params)
}

func replacePlaceholders(message string, params Params) string {
	if len(message) == 0 {
		return ""
	}
	for key, value := range params {
		message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
	}
	return message
}
