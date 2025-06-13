package group

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctrl := New()
	g := e.Group("/groups")
	g.POST("", ctrl.Create)
	g.DELETE("", ctrl.Delete)
	g.GET("/users", ctrl.ListUsers)
}
