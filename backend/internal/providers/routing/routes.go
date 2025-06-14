package routing

import (
	"github.com/labstack/echo/v4"
	ocservUserRoutes "ocserv-bakend/internal/api/ocserv_user"
	systemRoutes "ocserv-bakend/internal/api/system"
)

func Register(e *echo.Echo) {
	group := e.Group("/api")
	systemRoutes.Routes(group)
	ocservUserRoutes.Routes(group)
	//ocservUserRoutes.Routes(group)
}
