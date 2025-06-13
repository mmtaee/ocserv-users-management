package system

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/internal/repository"
	"ocserv-bakend/pkg/request"
)

type Controller struct {
	request    request.CustomRequestInterface
	systemRepo repository.SystemRepositoryInterface
}

func New() *Controller {
	return &Controller{
		request:    request.NewCustomRequest(),
		systemRepo: repository.NewSystemRepository(),
	}
}

// SystemInit
// @Summary      Get panel System init Config
// @Description  Get panel System init Config
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Failure      400 {object} request.ErrorResponse
// @Success      200  {object}  GetSystemInitResponse
// @Router       /system/Init [get]
func (ctrl *Controller) SystemInit(c echo.Context) error {
	config, err := ctrl.systemRepo.System(c.Request().Context())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusOK, nil)
		}
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, GetSystemInitResponse{
		GoogleCaptchaSiteKey: config.GoogleCaptchaSiteKey,
	})
}

// System
// @Summary      Get panel System Config
// @Description  Get panel System Config
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object}  GetSystemResponse
// @Router       /system [get]
func (ctrl *Controller) System(c echo.Context) error {
	config, err := ctrl.systemRepo.System(c.Request().Context())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusOK, nil)
		}
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, GetSystemResponse{
		GoogleCaptchaSiteKey:   config.GoogleCaptchaSiteKey,
		GoogleCaptchaSecretKey: config.GoogleCaptchaSecretKey,
	})
}

// SystemUpdate
// @Summary      Update panel System Config
// @Description  Update panel System Config
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Param        request    body  PatchSystemUpdateData   true "update config data"
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object}  GetSystemResponse
// @Router       /system [patch]
func (ctrl *Controller) SystemUpdate(c echo.Context) error {
	var data PatchSystemUpdateData
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	system := models.System{}

	if data.GoogleCaptchaSiteKey != nil {
		system.GoogleCaptchaSiteKey = *data.GoogleCaptchaSiteKey
	}
	if data.GoogleCaptchaSecretKey != nil {
		system.GoogleCaptchaSecretKey = *data.GoogleCaptchaSecretKey
	}

	updatedConfig, err := ctrl.systemRepo.SystemUpdate(c.Request().Context(), &system)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	return c.JSON(http.StatusOK, GetSystemResponse{
		GoogleCaptchaSiteKey:   updatedConfig.GoogleCaptchaSiteKey,
		GoogleCaptchaSecretKey: updatedConfig.GoogleCaptchaSecretKey,
	})
}
