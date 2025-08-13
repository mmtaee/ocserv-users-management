package log

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ocserv-bakend/internal/repository"
	"ocserv-bakend/pkg/request"
)

type Controller struct {
	request request.CustomRequestInterface
	logRepo repository.LogsRepositoryInterface
}

func New() *Controller {
	return &Controller{
		request: request.NewCustomRequest(),
		logRepo: repository.NewLogsRepository(),
	}
}

// UsersLogs 	 List of Users logs on self user model
//
// @Summary      List of Users logs on self user model
// @Description  List of Users logs on self user model
// @Tags         Logs(Users)
// @Accept       json
// @Produce      json
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 size query int false "Number of items per page" minimum(1) maximum(100) name(size)
// @Param 		 order query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200 {object} UsersLogsResponse
// @Router       /logs/users [get]
func (ctl *Controller) UsersLogs(c echo.Context) error {
	userUID := c.Get("userUID").(string)
	pagination := ctl.request.Pagination(c)

	logs, count, err := ctl.logRepo.UsersLogs(c.Request().Context(), pagination, userUID)

	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	return c.JSON(http.StatusOK, UsersLogsResponse{
		Meta: request.Meta{
			Page:         pagination.Page,
			PageSize:     pagination.PageSize,
			TotalRecords: count,
		},
		Result: logs,
	})
}

// AuditLogs 	 List of logs
//
// @Summary      List of logs
// @Description  List of logs
// @Tags         Logs
// @Accept       json
// @Produce      json
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 size query int false "Number of items per page" minimum(1) maximum(100) name(size)
// @Param 		 order query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Param 		 uid query string false "Search User by UID"
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Failure      403 {object} middlewares.PermissionDenied
// @Success      200 {object} UsersLogsResponse
// @Router       /logs/audit [get]
func (ctl *Controller) AuditLogs(c echo.Context) error {
	pagination := ctl.request.Pagination(c)
	uid := c.QueryParam("uid")
	logs, count, err := ctl.logRepo.Logs(c.Request().Context(), pagination, uid)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, UsersLogsResponse{
		Meta: request.Meta{
			Page:         pagination.Page,
			PageSize:     pagination.PageSize,
			TotalRecords: count,
		},
		Result: logs,
	})
}
