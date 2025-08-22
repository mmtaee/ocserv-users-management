package log

import (
	"api/pkg/routing/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	ctl := New()
	g := e.Group("/logs", middlewares.AuthMiddleware())

	g.GET("/users", ctl.UsersLogs)
	g.GET("/audit", ctl.AuditLogs)
}
