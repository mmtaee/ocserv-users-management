package system

import (
	"github.com/labstack/echo/v4"
	"ocserv-bakend/pkg/routing/middlewares"
)

func Routes(e *echo.Group) {
	ctrl := New()

	e.GET("/system/init", ctrl.SystemInit)
	e.POST("/system/users/login", ctrl.Login)

	g := e.Group("/system", middlewares.AuthMiddleware())
	g.GET("", ctrl.System)
	g.PATCH("", ctrl.SystemUpdate)

	g.POST("/system/users", ctrl.CreateUser, middlewares.AdminPermission())
	g.POST("/system/users/:uid/password", ctrl.ChangePassword, middlewares.AdminPermission())
	g.DELETE("/system/users/:uid", ctrl.DeleteUser, middlewares.AdminPermission())
}
