package jwe

import (
	"github.com/labstack/echo"
	"net/http"
)

type jwtExtractor func(echo.Context) (string, error)

var ErrJWTMissing = echo.NewHTTPError(http.StatusBadRequest, "missing or malformed jwt")

func JWTWithConfig(conf *Config) echo.MiddlewareFunc {
	extractor := jwtFromHeader(conf.GetBearer())
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, err := extractor(c)
			if err != nil {
				return err
			}
			token, err := conf.Decrypt(auth)
			if errV := token.Valid(); errV == nil && err == nil {
				// Store user information from token into context.
				c.Set("AccessToken", token)
				return next(c)

			}

			return &echo.HTTPError{
				Code:     http.StatusUnauthorized,
				Message:  "invalid or expired jwt",
				Internal: err,
			}
		}
	}
}

func jwtFromHeader(authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}
