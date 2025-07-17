package ocserv_user

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctl := New()
	g := e.Group("/ocserv/users")

	g.GET("", ctl.OcservUsers)
	g.POST("", ctl.CreateOcservUser)
	g.PATCH("/:uid", ctl.UpdateOcservUser)
	g.DELETE("/:uid", ctl.DeleteOcservUser)
	g.POST("/:uid/lock", ctl.LockOcservUser)
	g.POST("/:uid/unlock", ctl.UnLockOcservUser)
	g.POST("/:username/disconnect", ctl.DisconnectOcservUser)
	g.GET("/:uid/statistics", ctl.StatisticsOcservUser)
}
