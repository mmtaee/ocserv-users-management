package system

import (
	"github.com/labstack/echo/v4"
	"ocserv-bakend/pkg/routing/middlewares"
)

func Routes(e *echo.Group) {
	ctrl := New()

	e.GET("/system/init", ctrl.SystemInit)
	e.POST("/system/setup", ctrl.SetupSystem)
	e.POST("/system/users/login", ctrl.Login)

	g := e.Group("/system", middlewares.AuthMiddleware())
	g.GET("", ctrl.System)
	g.PATCH("", ctrl.SystemUpdate)

	g.POST("/users", ctrl.CreateUser, middlewares.AdminPermission())
	g.POST("/users/:uid/password", ctrl.ChangeUserPasswordByAdmin, middlewares.AdminPermission())
	g.DELETE("/users/:uid", ctrl.DeleteUser, middlewares.AdminPermission())
	g.POST("/users/password", ctrl.ChangePasswordBySelf)
	g.GET("/users/profile", ctrl.Profile)
}
