package ocserv_user

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctrl := New()
	g := e.Group("/ocserv/users")

	g.GET("", ctrl.OcservUsers)
	g.POST("", ctrl.CreateOcservUser)
	g.PATCH("/:uid", ctrl.UpdateOcservUser)
	g.DELETE("/:uid", ctrl.DeleteOcservUser)
	g.POST("/:uid/lock", ctrl.LockOcservUser)
	g.POST("/:uid/unlock", ctrl.UnLockOcservUser)
}
