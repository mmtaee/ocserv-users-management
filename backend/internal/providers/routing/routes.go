package routing

import (
	"github.com/labstack/echo/v4"
	homeRoutes "ocserv-bakend/internal/api/home"
	logRoutes "ocserv-bakend/internal/api/log"
	occtlRoutes "ocserv-bakend/internal/api/occtl"
	ocservGroupRoutes "ocserv-bakend/internal/api/ocserv_group"
	ocservUserRoutes "ocserv-bakend/internal/api/ocserv_user"
	systemRoutes "ocserv-bakend/internal/api/system"
)

func Register(e *echo.Echo) {
	group := e.Group("/api")
	systemRoutes.Routes(group)
	ocservUserRoutes.Routes(group)
	ocservGroupRoutes.Routes(group)
	homeRoutes.Routes(group)
	occtlRoutes.Routes(group)
	logRoutes.Routes(group)
}
