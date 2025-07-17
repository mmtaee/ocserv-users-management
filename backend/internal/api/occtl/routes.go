package occtl

import (
	"github.com/labstack/echo/v4"
	"ocserv-bakend/pkg/routing/middlewares"
)

func Routes(e *echo.Group) {
	ctl := New()
	g := e.Group("/occtl")
	g.GET("/server_info", ctl.ServerInfo)
	g.GET("/commands", ctl.Commands, middlewares.AuthMiddleware())
}
