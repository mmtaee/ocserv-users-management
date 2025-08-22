package occtl

import (
	"api/pkg/routing/middlewares"
	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Group) {
	ctl := New()
	g := e.Group("/occtl")
	g.GET("/server_info", ctl.ServerInfo)
	g.GET("/commands", ctl.Commands, middlewares.AuthMiddleware())
}
