package log

import (
	"github.com/labstack/echo/v4"
	"ocserv-bakend/pkg/routing/middlewares"
)

func Routes(e *echo.Group) {
	ctl := New()
	g := e.Group("/logs", middlewares.AuthMiddleware())

	g.GET("/users", ctl.UsersLogs)
	g.GET("/audit", ctl.AuditLogs)
}
