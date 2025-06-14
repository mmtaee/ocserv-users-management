package ocserv_group

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctrl := New()
	g := e.Group("/ocserv/groups")
	g.GET("", ctrl.OcservGroups)

}
