package system

import (
	"github.com/labstack/echo/v4"
	"ocserv-bakend/pkg/routing/middlewares"
)

func Routes(e *echo.Group) {
	ctrl := New()

	e.GET("/system/init", ctrl.SystemInit)

	g := e.Group("/system", middlewares.AuthMiddleware())
	g.GET("", ctrl.System)
	g.PATCH("", ctrl.SystemUpdate)

}
