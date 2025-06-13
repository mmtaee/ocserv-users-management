package system

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/internal/repository"
	"ocserv-bakend/pkg/crypto"
	"ocserv-bakend/pkg/request"
	"ocserv-bakend/pkg/utils/captcha"
)

type Controller struct {
	request         request.CustomRequestInterface
	systemRepo      repository.SystemRepositoryInterface
	userRepo        repository.UserRepositoryInterface
	captchaVerifier captcha.GoogleCaptchaInterface
	cryptoRepo      crypto.CustomPasswordInterface
}

func New() *Controller {
	return &Controller{
		request:         request.NewCustomRequest(),
		systemRepo:      repository.NewSystemRepository(),
		userRepo:        repository.NewUserRepository(),
		captchaVerifier: captcha.NewGoogleVerifier(),
		cryptoRepo:      crypto.NewCustomPassword(),
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
// @Param        request    body  PatchSystemUpdateData   true "update system config data"
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

// Login		 Admin users login
//
// @Summary      Admin users login
// @Description  Admin users login with Google captcha(captcha site key required in get config api)
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Param        request    body  LoginData   true "login data"
// @Failure      400 {object} request.ErrorResponse
// @Success      201  {object}  UserLoginResponse
// @Router       /users/login [post]
func (ctrl *Controller) Login(c echo.Context) error {
	var data LoginData
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	system, err := ctrl.systemRepo.System(c.Request().Context())
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	if secretKey := system.GoogleCaptchaSecretKey; secretKey != "" {
		ctrl.captchaVerifier.SetSecretKey(secretKey)
		ctrl.captchaVerifier.Verify(data.Token)
		if !ctrl.captchaVerifier.IsValid() {
			return ctrl.request.BadRequest(c, errors.New("captcha challenge failed"))
		}
	}

	user, err := ctrl.userRepo.GetByUsername(c.Request().Context(), data.Username)
	if err != nil {
		return ctrl.request.BadRequest(c, errors.New("invalid username or password"))
	}

	token, err := ctrl.userRepo.CreateToken(c.Request().Context(), user.ID, user.UID, true, user.IsAdmin)
	if err != nil {
		return ctrl.request.BadRequest(c, err, "user created")
	}

	return c.JSON(http.StatusOK, UserLoginResponse{
		User:  user,
		Token: token,
	})
}

// CreateUser	 Create user
//
// @Summary      Create user
// @Description  Create user Admin or simple
// @Tags         Panel
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request    body  CreateUserData   true "create user data"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Failure      403 {object} middlewares.PermissionDenied
// @Success      201  {object}  models.User
// @Router       /users/login [post]
func (ctrl *Controller) CreateUser(c echo.Context) error {
	var data CreateUserData
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	passwd := ctrl.cryptoRepo.CreatePassword(data.Password)

	user := &models.User{
		Username: data.Username,
		Password: passwd.Hash,
		Salt:     passwd.Salt,
		IsAdmin:  data.Admin,
	}
	newUser, err := ctrl.userRepo.CreateUser(c.Request().Context(), user)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusCreated, newUser)
}
