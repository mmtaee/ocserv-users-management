package ocserv_group

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"ocserv-bakend/internal/models"
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

// OcservGroupsLookup 	 List of Ocserv group names
//
// @Summary      List of Ocserv group names
// @Description  List of Ocserv group names
// @Tags         Ocserv(Groups)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200 {array}  string
// @Router       /ocserv/groups/lookup [get]
func (ctrl *Controller) OcservGroupsLookup(c echo.Context) error {
	groups, err := ctrl.ocservGroupRepo.GroupsLookup(c.Request().Context())
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	groups = append([]string{"defaults"}, groups...)
	return c.JSON(http.StatusOK, groups)
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

// CreateOcservGroup 	     Ocserv Group creation
//
// @Summary      Ocserv Group creation
// @Description  Ocserv Group creation
// @Tags         Ocserv(Groups)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request    body  CreateOcservGroupData  true "ocserv group create data"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      201  {object} models.OcservGroup
// @Router       /ocserv/groups [post]
func (ctrl *Controller) CreateOcservGroup(c echo.Context) error {
	var data CreateOcservGroupData
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	ocservGroup := models.OcservGroup{
		Name:   data.Name,
		Config: data.Config,
	}

	newOcservGroup, err := ctrl.ocservGroupRepo.Create(c.Request().Context(), &ocservGroup)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, newOcservGroup)
}

// UpdateOcservGroup 	     Ocserv Group update
//
// @Summary      Ocserv Group update
// @Description  Ocserv Group update
// @Tags         Ocserv(Groups)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path int true "Ocserv Group ID"
// @Param        request    body  UpdateOcservGroupData  true "ocserv group create data"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      201  {object} models.OcservGroup
// @Router       /ocserv/groups/{uid} [patch]
func (ctrl *Controller) UpdateOcservGroup(c echo.Context) error {
	groupUID := c.Param("uid")
	if groupUID == "" {
		return ctrl.request.BadRequest(c, errors.New("group uid is empty"))
	}

	var data UpdateOcservGroupData
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	ocservGroup, err := ctrl.ocservGroupRepo.GetByUID(c.Request().Context(), groupUID)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	ocservGroup.Config = data.Config
	updatedOcservGroup, err := ctrl.ocservGroupRepo.Update(c.Request().Context(), ocservGroup)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, updatedOcservGroup)
}

// DeleteOcservGroup 	     Ocserv Group delete
//
// @Summary      Ocserv Group delete
// @Description  Ocserv Group delete
// @Tags         Ocserv(Groups)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path int true "Ocserv Group ID"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      204  {object} nil
// @Router       /ocserv/groups/{uid} [delete]
func (ctrl *Controller) DeleteOcservGroup(c echo.Context) error {
	groupUID := c.Param("uid")
	if groupUID == "" {
		return ctrl.request.BadRequest(c, errors.New("group uid is empty"))
	}

	err := ctrl.ocservGroupRepo.Delete(c.Request().Context(), groupUID)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}
