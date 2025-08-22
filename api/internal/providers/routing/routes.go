package routing

import (
	homeRoutes "api/internal/api/home"
	logRoutes "api/internal/api/log"
	occtlRoutes "api/internal/api/occtl"
	ocservGroupRoutes "api/internal/api/ocserv_group"
	ocservUserRoutes "api/internal/api/ocserv_user"
	systemRoutes "api/internal/api/system"
	"github.com/labstack/echo/v4"
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
