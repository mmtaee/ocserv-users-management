package routing

import (
	"github.com/labstack/echo/v4"
	systemRoutes "ocserv-bakend/internal/api/system"
)

func Register(e *echo.Echo) {
	group := e.Group("/api")
	systemRoutes.Routes(group)
	//userRoutes.Routes(group)
	//ocservUserRoutes.Routes(group)
}
