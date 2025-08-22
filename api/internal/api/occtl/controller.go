package occtl

import (
	"api/internal/repository"
	"api/pkg/config"
	ocApi "api/pkg/oc_api"
	"api/pkg/request"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strings"
)

type Controller struct {
	request   request.CustomRequestInterface
	occtlRepo repository.OcctlRepositoryInterface
	ocApi     ocApi.OcOcctlApiRepositoryInterface
}

func New() *Controller {
	webhookApi := config.Get().WebhookApi
	return &Controller{
		request:   request.NewCustomRequest(),
		occtlRepo: repository.NewOcctlRepository(),
		ocApi:     ocApi.NewOcctlApiRepository(webhookApi),
	}
}

// ServerInfo 	 Server information
//
// @Summary      Server information
// @Description  Server information
// @Tags         OCCTL
// @Accept       json
// @Produce      json
// @Failure      400 {object} request.ErrorResponse
// @Success      200  {object}  models.ServerVersion
// @Router       /occtl/server_info [get]
func (ctl *Controller) ServerInfo(c echo.Context) error {
	version, err := ctl.occtlRepo.Version(c.Request().Context())
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, version)
}

// Commands 	 Occtl Commands
//
// @Summary      Occtl Commands
// @Description  Occtl Commands
// @Tags         OCCTL
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        action  query   int     true   "Command Action ID (1 to 15)"
// @Param        value   query   string  false  "Optional parameter depending on command"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object}  string
// @Router       /occtl/commands [get]
func (ctl *Controller) Commands(c echo.Context) error {
	var data CommandParamsData
	if err := c.Bind(&data); err != nil {
		return ctl.request.BadRequest(c, err)
	}
	var (
		results []byte
		err     error
	)
	ctx := c.Request().Context()
	actions := map[int]func(context.Context, string) ([]byte, error){
		1: func(ctx context.Context, _ string) ([]byte, error) { return ctl.ocApi.OnlineUsersInfo(ctx) },
		2: ctl.ocApi.ShowUserByUsername,
		3: ctl.ocApi.ShowUserByID,
		4: ctl.ocApi.Disconnect,
		5: func(ctx context.Context, _ string) ([]byte, error) { return ctl.ocApi.ShowSessionsAll(ctx) },
		6: func(ctx context.Context, _ string) ([]byte, error) { return ctl.ocApi.ShowSessionsAllValid(ctx) },
		7: ctl.ocApi.ShowSessionBySID,
		8: func(ctx context.Context, _ string) ([]byte, error) { return ctl.ocApi.IPBans(ctx) },
		9: ctl.ocApi.UnbanIP,
		10: func(ctx context.Context, _ string) ([]byte, error) {
			ctx = context.WithValue(ctx, "format", "json")
			return ctl.ocApi.Status(ctx)
		},
		11: func(ctx context.Context, _ string) ([]byte, error) { return ctl.ocApi.ShowEvent(ctx) },
		12: func(ctx context.Context, _ string) ([]byte, error) { return ctl.ocApi.IRoutes(ctx) },
	}

	if data.Action == 13 {
		err = ctl.ocApi.Reload(ctx)
		results = []byte(`{"message": "Server reload has been scheduled successfully"}`)
	} else if handler, exists := actions[data.Action]; exists {
		results, err = handler(ctx, data.Value)
	} else {
		return ctl.request.BadRequest(c, fmt.Errorf("unknown action %d", data.Action))
	}

	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	log.Println(string(results))
	return c.JSON(http.StatusOK, strings.TrimSpace(string(results)))
}
