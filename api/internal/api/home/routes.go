package home

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctl := New()
	e.GET("/home", ctl.Home)
}
