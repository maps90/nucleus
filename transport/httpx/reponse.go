package httpx

import (
	"net/http"

	"github.com/labstack/echo"
)

// JSONResponse struct
type JSONResponse struct {
	RequestID string      `json:"request_id"`
	Code      int         `json:"status_code"`
	Result    interface{} `json:"data,omitempty"`
}

// JSONR return JSON response
func (c *Context) JSONR(data interface{}) error {
	response := &JSONResponse{
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
		Code:      http.StatusOK,
		Result:    data,
	}

	return c.JSON(http.StatusOK, response)
}
