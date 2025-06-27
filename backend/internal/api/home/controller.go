package home

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ocserv-bakend/internal/repository"
	"ocserv-bakend/pkg/request"
)

type Controller struct {
	request        request.CustomRequestInterface
	occtlRepo      repository.OcctlRepositoryInterface
	ocservUserRepo repository.OcservUserRepositoryInterface
}

func New() *Controller {
	return &Controller{
		request:        request.NewCustomRequest(),
		occtlRepo:      repository.NewOcctlRepository(),
		ocservUserRepo: repository.NewtOcservUserRepository(),
	}
}

// Home 	     Content of home
//
// @Summary      Content of home
// @Description  Content of home
// @Tags         Home
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object}  HomeResponse
// @Router       /home [get]
func (ctrl *Controller) Home(c echo.Context) error {
	ctx := c.Request().Context()

	serverStatus, err := ctrl.occtlRepo.Stats(ctx)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	status := SplitStatsText(serverStatus)
	stats, err := ctrl.ocservUserRepo.TenDaysStats(ctx)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	onlineUsers, err := ctrl.occtlRepo.OnlineUsers(ctx)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	
	return c.JSON(http.StatusOK, HomeResponse{
		Status:     status,
		Stats:      stats,
		OnlineUser: onlineUsers,
	})
}
