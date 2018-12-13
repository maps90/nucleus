package errors

type APIError struct {
	RequestID string `json:"request_id"`
	// Status is http status code
	Status int `json:"status_code"`
	// ErrorCode is the code uniquely identifying an error
	ErrorCode string `json:"error_code"`
	// Message is the error message that may be displayed to end users
	Message string `json:"error_message"`
	// DeveloperMessage is the error message that is mainly meant for developers
	DeveloperMessage string `json:"developer_message,omitempty"`
	// Details specifies the additional error information
	Details interface{} `json:"details,omitempty"`
}

// Error returns the error message.
func (e APIError) Error() string {
	return e.Message
}

// StatusCode return http status code
func (e APIError) StatusCode() int {
	return e.Status
}
