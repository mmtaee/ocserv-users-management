package occtl

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctrl := New()
	e.GET("/version", ctrl.Version)

	g := e.Group("/occtl")
	g.GET("/online-users/info", ctrl.OnlineUsersInfo)

	g.GET("/online-users", ctrl.OnlineUsers)
	g.POST("/disconnect/:username", ctrl.DisconnectUser)
	g.POST("/reload", ctrl.Reload)
	g.GET("/ip-bans", ctrl.ShowIPBans)
	g.DELETE("/unban-ip/:ip", ctrl.UnbanIP)
	g.GET("/status", ctrl.ShowStatus)
	g.GET("/iroutes", ctrl.ShowIRoutes)
	g.GET("/users/id/:id", ctrl.ShowUserByID)
	g.GET("/users/:username", ctrl.ShowUser)
	g.GET("/sessions/sid/:sid", ctrl.ShowSession)
	g.GET("/sessions/valid", ctrl.ShowSessionsValid)
	g.GET("/sessions", ctrl.ShowSessionsALL)
	g.GET("/events", ctrl.ShowEvent)
}
