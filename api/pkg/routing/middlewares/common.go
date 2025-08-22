package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func UnauthorizedError(c echo.Context, msg string) error {
	// "invalid authentication credentials"
	return c.JSON(http.StatusUnauthorized, Unauthorized{Error: msg})
}

func PermissionDeniedError(c echo.Context, msg string) error {
	return c.JSON(http.StatusForbidden, PermissionDenied{Error: msg})
}
