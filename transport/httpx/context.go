package httpx

import (
	"context"

	"github.com/labstack/echo"
)

type (
	// Context struct
	Context struct {
		echo.Context
	}

	// ContextFunc typefunc
	ContextFunc func(*Context) error

	key string
)

// KeyHandler custom handler
var KeyHandler key = "custom_handler"

// NewHandler generate a base handler
func NewHandler(ctxFunc ContextFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := context.WithValue(c.Request().Context(), KeyHandler, nil)
		c.SetRequest(c.Request().WithContext(ctx))

		return ctxFunc(&Context{c})
	}
}
