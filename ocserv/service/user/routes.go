package user

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctrl := New()
	g := e.Group("/users")
	g.POST("", ctrl.Create)
	g.POST("/:username/lock", ctrl.Lock)
	g.POST("/:username/unlock", ctrl.Unlock)
	g.DELETE("/:username", ctrl.Delete)
	g.POST("/:username/config", ctrl.CreateConfig)
	g.DELETE("/:username/config", ctrl.DeleteConfig)
}
