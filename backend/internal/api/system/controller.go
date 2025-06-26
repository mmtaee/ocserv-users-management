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
	"strings"
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

// SetupSystem
// @Summary      Setup user and system config
// @Description  Setup user and system config
// @Tags         System
// @Accept       json
// @Produce      json
// @Param        request  body  SetupSystem   true "system setup data"
// @Failure      400 {object} request.ErrorResponse
// @Success      201  {object}  SetupSystemResponse
// @Router       /system/setup [post]
func (ctrl *Controller) SetupSystem(c echo.Context) error {
	if _, err := ctrl.systemRepo.System(c.Request().Context()); err == nil {
		return ctrl.request.BadRequest(c, errors.New("the system is already configured"))
	}

	var data SetupSystem
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	passwd := ctrl.cryptoRepo.CreatePassword(data.Password)

	user := &models.User{
		Username: strings.ToLower(data.Username),
		Password: passwd.Hash,
		Salt:     passwd.Salt,
		IsAdmin:  true,
	}

	system := &models.System{
		GoogleCaptchaSiteKey:   data.GoogleCaptchaSiteKey,
		GoogleCaptchaSecretKey: data.GoogleCaptchaSecretKey,
	}
	newUser, newSystem, err := ctrl.systemRepo.SystemSetup(c.Request().Context(), user, system)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	token, err := ctrl.userRepo.CreateToken(c.Request().Context(), newUser.ID, newUser.UID, true, newUser.IsAdmin)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	return c.JSON(
		http.StatusCreated,
		SetupSystemResponse{
			User:   *newUser,
			System: *newSystem,
			Token:  token,
		},
	)
}

// SystemInit
// @Summary      Get panel System init Config
// @Description  Get panel System init Config
// @Tags         System
// @Accept       json
// @Produce      json
// @Failure      400 {object} request.ErrorResponse
// @Success      200  {object}  GetSystemInitResponse
// @Router       /system/init [get]
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
// @Tags         System
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
// @Tags         System
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
// @Tags         System(Users)
// @Accept       json
// @Produce      json
// @Param        request    body  LoginData   true "login data"
// @Failure      400 {object} request.ErrorResponse
// @Success      201  {object}  UserLoginResponse
// @Router       /system/users/login [post]
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
// @Tags         System(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request    body  CreateUserData   true "create user data"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Failure      403 {object} middlewares.PermissionDenied
// @Success      201  {object}  models.User
// @Router       /system/users/login [post]
func (ctrl *Controller) CreateUser(c echo.Context) error {
	var data CreateUserData
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	passwd := ctrl.cryptoRepo.CreatePassword(data.Password)

	user := &models.User{
		Username: strings.ToLower(data.Username),
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

// Users 		 List of Users
//
// @Summary      List of Admin or simple users
// @Description  List of Admin or simple users
// @Tags         System(Users)
// @Accept       json
// @Produce      json
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 page_size query int false "Number of items per page" minimum(1) maximum(100)
// @Param 		 order query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Failure      403 {object} middlewares.PermissionDenied
// @Success      200  {object}  UsersResponse
// @Router       /system/users [get]
func (ctrl *Controller) Users(c echo.Context) error {
	pagination := ctrl.request.Pagination(c)

	users, total, err := ctrl.userRepo.Users(c.Request().Context(), pagination)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}

	return c.JSON(http.StatusOK, UsersResponse{
		Meta: request.Meta{
			Page:         pagination.Page,
			PageSize:     pagination.PageSize,
			TotalRecords: total,
		},
		Result: users,
	})
}

// ChangeUserPasswordByAdmin 		 Change user password by admin
//
// @Summary      Change user password by admin
// @Description  Change user password by admin
// @Tags         System(Users)
// @Accept       json
// @Produce      json
// @Param 		 uid path int true "User ID"
// @Param        request    body  ChangeUserPassword  true "user new password"
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Failure      403 {object} middlewares.PermissionDenied
// @Success      200  {object}  UsersResponse
// @Router       /system/users/{uid}/password [post]
func (ctrl *Controller) ChangeUserPasswordByAdmin(c echo.Context) error {
	userID := c.Param("uid")

	var data ChangeUserPassword
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	passwd := ctrl.cryptoRepo.CreatePassword(data.Password)

	err := ctrl.userRepo.ChangePassword(c.Request().Context(), userID, passwd.Hash, passwd.Salt)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// DeleteUser 	 Delete simple user
//
// @Summary      Delete simple user
// @Description  Delete simple user
// @Tags         System(Users)
// @Accept       json
// @Produce      json
// @Param 		 uid path int true "User ID"
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Failure      403 {object} middlewares.PermissionDenied
// @Success      204  {object}  nil
// @Router       /system/users/{uid} [delete]
func (ctrl *Controller) DeleteUser(c echo.Context) error {
	userID := c.Param("uid")
	err := ctrl.userRepo.DeleteUser(c.Request().Context(), userID)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

// ChangePasswordBySelf 		 Change user password by self
//
// @Summary      Change user password by self
// @Description  Change user password by self
// @Tags         System(Users)
// @Accept       json
// @Produce      json
// @Param        request body  ChangeUserPassword  true "user new password"
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object}  UsersResponse
// @Router       /system/users/password [post]
func (ctrl *Controller) ChangePasswordBySelf(c echo.Context) error {
	userID := c.Get("userUID").(string)

	var data ChangeUserPassword
	if err := ctrl.request.DoValidate(c, &data); err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	passwd := ctrl.cryptoRepo.CreatePassword(data.Password)
	err := ctrl.userRepo.ChangePassword(c.Request().Context(), userID, passwd.Hash, passwd.Salt)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// Profile 		 Get User Profile
//
// @Summary      Get User Profile
// @Description  Get User Profile
// @Tags         System(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object}  models.User
// @Router       /system/users/profile [get]
func (ctrl *Controller) Profile(c echo.Context) error {
	userUID := c.Get("userUID").(string)
	user, err := ctrl.userRepo.GetByUID(c.Request().Context(), userUID)
	if err != nil {
		return ctrl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, user)
}
