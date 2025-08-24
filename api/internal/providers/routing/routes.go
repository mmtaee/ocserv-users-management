package routing

import (
	//homeRoutes "api/internal/api/home"
	//logRoutes "api/internal/api/log"
	//occtlRoutes "api/internal/api/occtl"
	//ocservGroupRoutes "api/internal/api/ocserv_group"
	//ocservUserRoutes "api/internal/api/ocserv_user"
	"github.com/labstack/echo/v4"
	ocservGroupRoutes "github.com/mmtaee/ocserv-users-management/api/internal/services/ocserv_group"
	ocservUserRoutes "github.com/mmtaee/ocserv-users-management/api/internal/services/ocserv_user"
	systemRoutes "github.com/mmtaee/ocserv-users-management/api/internal/services/system"
)

func Register(e *echo.Echo) {
	group := e.Group("/api")
	systemRoutes.Routes(group)
	ocservGroupRoutes.Routes(group)
	ocservUserRoutes.Routes(group)
	//homeRoutes.Routes(group)
	//occtlRoutes.Routes(group)
	//logRoutes.Routes(group)
}
