package ocserv_group

import "github.com/labstack/echo/v4"

func Routes(e *echo.Group) {
	ctl := New()
	g := e.Group("/ocserv/groups")
	g.GET("", ctl.OcservGroups)
	g.GET("/lookup", ctl.OcservGroupsLookup)
	g.POST("", ctl.CreateOcservGroup)
	g.PATCH("/:id", ctl.UpdateOcservGroup)
	g.DELETE("/:id", ctl.DeleteOcservGroup)
	g.GET("/defaults", ctl.GetDefaultsGroup)
	g.PATCH("/defaults", ctl.UpdateDefaultsGroup)
}
