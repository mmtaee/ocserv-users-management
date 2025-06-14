package ocserv_group

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ocserv-bakend/internal/repository"
	"ocserv-bakend/pkg/request"
)

type Controller struct {
	request         request.CustomRequestInterface
	ocservGroupRepo repository.OcservGroupRepositoryInterface
}

func New() *Controller {
	return &Controller{
		request:         request.NewCustomRequest(),
		ocservGroupRepo: repository.NewOcservGroupRepository(),
	}
}

// OcservGroups 	 List of Ocserv groups
//
// @Summary      List of Ocserv groups
// @Description  List of Ocserv groups
// @Tags         Ocserv(Groups)
// @Accept       json
// @Produce      json
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 page_size query int false "Number of items per page" minimum(1) maximum(100)
// @Param 		 order query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object}  OcservGroupsResponse
// @Router       /ocserv/groups [get]
func (ctrl *Controller) OcservGroups(c echo.Context) error {
	pagination := ctrl.request.Pagination(c)

	ocservGroup, total, err := ctrl.ocservGroupRepo.Groups(c.Request().Context(), pagination)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	return c.JSON(http.StatusOK, OcservGroupsResponse{
		Meta: request.Meta{
			Page:         pagination.Page,
			PageSize:     pagination.PageSize,
			TotalRecords: total,
		},
		Result: ocservGroup,
	})
}
