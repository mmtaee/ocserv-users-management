package home

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctrl := New()
	e.GET("/home", ctrl.Home)
}
