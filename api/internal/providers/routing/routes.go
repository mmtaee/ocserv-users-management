package routing

import (
	"github.com/labstack/echo/v4"
	occtlRoutes "github.com/mmtaee/ocserv-users-management/api/internal/services/occtl"
	ocservGroupRoutes "github.com/mmtaee/ocserv-users-management/api/internal/services/ocserv_group"
	ocservUserRoutes "github.com/mmtaee/ocserv-users-management/api/internal/services/ocserv_user"
	systemRoutes "github.com/mmtaee/ocserv-users-management/api/internal/services/system"
)

func Register(e *echo.Echo) {
	group := e.Group("/api")
	systemRoutes.Routes(group)
	ocservGroupRoutes.Routes(group)
	ocservUserRoutes.Routes(group)
	occtlRoutes.Routes(group)
	//homeRoutes.Routes(group)
	//logRoutes.Routes(group)
}
