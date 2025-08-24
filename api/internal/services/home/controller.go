package home

import (
	"github.com/labstack/echo/v4"
	"github.com/mmtaee/ocserv-users-management/api/internal/repository"
	"github.com/mmtaee/ocserv-users-management/api/pkg/request"
	"github.com/mmtaee/ocserv-users-management/common/models"
	"log"
	"net/http"
	"sync"
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
// @Success      200  {object}  GetHomeResponse
// @Router       /home [get]
func (ctl *Controller) Home(c echo.Context) error {
	ctx := c.Request().Context()

	var (
		status      ServerStatusResponse
		statistics  *[]models.DailyTraffic
		onlineUsers *[]models.OnlineUserSession
		ipBans      *[]models.IPBanPoints
		errs        = make(chan error, 4)
		wg          sync.WaitGroup
	)

	wg.Add(4)

	go func() {
		defer wg.Done()
		serverStatus, err := ctl.occtlRepo.Status()
		if err != nil {
			errs <- err
			return
		}
		if serverStatusMap, ok := serverStatus.(map[string]interface{}); ok {
			status = ParseServerStatus(serverStatusMap)
		}
	}()

	go func() {
		defer wg.Done()
		data, err := ctl.ocservUserRepo.TenDaysStats(ctx)
		if err != nil {
			errs <- err
			return
		}
		statistics = data
	}()

	go func() {
		defer wg.Done()
		users, err := ctl.occtlRepo.OnlineUsersInfo()
		if err != nil {
			errs <- err
			return
		}
		onlineUsers = users
	}()

	go func() {
		defer wg.Done()
		ips, err := ctl.occtlRepo.IPBans()
		if err != nil {
			errs <- err
			return
		}
		ipBans = ips
	}()

	wg.Wait()
	close(errs)

	if err := <-errs; err != nil {
		log.Println("error in Home handler:", err)
		return ctl.request.BadRequest(c, err)
	}

	resp := GetHomeResponse{
		ServerStatus: status,
		Statistics:   statistics,
		OnlineUser:   onlineUsers,
		IPBans:       ipBans,
	}

	return c.JSON(http.StatusOK, resp)
}
