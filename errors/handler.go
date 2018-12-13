package errors

import (
	"database/sql"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo"
)

// ErrorHandler override echo.HTTPErrorHandler
func ErrorHandler(err error, c echo.Context) {
	var e = &Errors{c}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, e.NotFound("the requested resource"))
		return
	}

	switch err.(type) {
	case *APIError:
		c.JSON(err.(*APIError).StatusCode(), err)
	case validation.Errors:
		c.JSON(http.StatusBadRequest, e.InvalidData(err.(validation.Errors)))
	case *echo.HTTPError:
		switch err.(*echo.HTTPError).Code {
		case http.StatusUnauthorized:
			c.JSON(http.StatusUnauthorized, e.Unauthorized(err.Error()))
		case http.StatusNotFound, http.StatusMethodNotAllowed:
			c.JSON(http.StatusNotFound, e.NotFound("the requested resource"))
		default:
			errEcho := err.(*echo.HTTPError)
			c.JSON(errEcho.Code, NewAPIError(c, http.StatusText(errEcho.Code), Params{"error": errEcho.Message}))
		}
	default:
		c.JSON(http.StatusInternalServerError, e.InternalServerError(err))
	}
}
