package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Middleware is a simple middleware that checks the request for authentication
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// obviously this is not real authentication and is just illustrative of what you can do here
			username := c.Request().URL.Query().Get("username")
			if username == "reject" {
				return echo.ErrUnauthorized.SetInternal(errors.Errorf("user %s is not authorized", username))
			}

			return next(c)
		}
	}
}
