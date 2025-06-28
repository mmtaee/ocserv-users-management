package occtl

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctrl := New()
	g := e.Group("/occtl")
	g.GET("/online-users", ctrl.OnlineUsers)
	g.GET("/online-users/info", ctrl.OnlineUsersInfo)
	g.POST("/disconnect/:username", ctrl.DisconnectUser)
	g.POST("/reload", ctrl.Reload)
	g.GET("/ip-bans", ctrl.ShowIPBans)
	g.DELETE("/ip-bans/:ip", ctrl.UnbanIP)
	g.GET("/status", ctrl.ShowStatus)
	g.GET("/iroutes", ctrl.ShowIRoutes)
	g.GET("/user/:username", ctrl.ShowUser)
}
