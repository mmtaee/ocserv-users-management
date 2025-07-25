package ocserv_user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/internal/repository"
	"ocserv-bakend/pkg/config"
	ocApi "ocserv-bakend/pkg/oc_api"
	"ocserv-bakend/pkg/request"
	"slices"
	"time"
)

type Controller struct {
	request        request.CustomRequestInterface
	ocservUserRepo repository.OcservUserRepositoryInterface
	ocRepo         ocApi.OcOcctlApiRepositoryInterface
}

func New() *Controller {
	apiURLService := config.Get().APIURLService
	return &Controller{
		request:        request.NewCustomRequest(),
		ocservUserRepo: repository.NewtOcservUserRepository(),
		ocRepo:         ocApi.NewOcctlApiRepository(apiURLService),
	}
}

// OcservUsers 	 List of Ocserv Users
//
// @Summary      List of Ocserv Users
// @Description  List of Ocserv Users
// @Tags         Ocserv(Users)
// @Accept       json
// @Produce      json
// @Param 		 page query int false "Page number, starting from 1" minimum(1)
// @Param 		 size query int false "Number of items per page" minimum(1) maximum(100) name(size)
// @Param 		 order query string false "Field to order by"
// @Param 		 sort query string false "Sort order, either ASC or DESC" Enums(ASC, DESC)
// @Param        Authorization header string true "Bearer TOKEN"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object}  OcservUsersResponse
// @Router       /ocserv/users [get]
func (ctl *Controller) OcservUsers(c echo.Context) error {
	pagination := ctl.request.Pagination(c)

	ocservUser, total, err := ctl.ocservUserRepo.Users(c.Request().Context(), pagination)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	onlineUserBytes, err := ctl.ocRepo.OnlineUsers(c.Request().Context())
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	var onlineUsers struct {
		Users []string `json:"users"`
	}

	err = json.Unmarshal(onlineUserBytes, &onlineUsers)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	for i := range *ocservUser {
		user := &(*ocservUser)[i]
		if slices.Contains(onlineUsers.Users, user.Username) {
			user.IsOnline = true
		}
	}

	return c.JSON(http.StatusOK, OcservUsersResponse{
		Meta: request.Meta{
			Page:         pagination.Page,
			TotalRecords: total,
			PageSize:     pagination.PageSize,
		},
		Result: ocservUser,
	})
}

// CreateOcservUser 	     Ocserv User creation
//
// @Summary      Ocserv User creation
// @Description  Ocserv User creation
// @Tags         Ocserv(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param        request    body  CreateOcservUserData  true "ocserv user create data"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      201  {object} models.OcservUser
// @Router       /ocserv/users [post]
func (ctl *Controller) CreateOcservUser(c echo.Context) error {
	var data CreateOcservUserData
	if err := ctl.request.DoValidate(c, &data); err != nil {
		return ctl.request.BadRequest(c, err)
	}

	expireAt, err := time.Parse("2006-01-02", data.ExpireAt)
	if err != nil {
		expireAt, _ = time.Parse(
			"2006-01-02",
			time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
		)
	}

	if data.TrafficType == models.Free {
		data.TrafficSize = 0
	}

	ocUser := &models.OcservUser{
		Username:    data.Username,
		Password:    data.Password,
		Group:       data.Group,
		ExpireAt:    &expireAt,
		TrafficSize: data.TrafficSize,
		TrafficType: data.TrafficType,
		Config:      data.Config,
	}

	user, err := ctl.ocservUserRepo.Create(c.Request().Context(), ocUser)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateOcservUser 	     Ocserv User update
//
// @Summary      Ocserv User update
// @Description  Ocserv User update
// @Tags         Ocserv(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Param        request    body  UpdateOcservUserData  true "ocserv user update data"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      201  {object} models.OcservUser
// @Router       /ocserv/users/{uid} [patch]
func (ctl *Controller) UpdateOcservUser(c echo.Context) error {
	userID := c.Param("uid")
	if userID == "" {
		return ctl.request.BadRequest(c, errors.New("user id is required"))
	}

	var data UpdateOcservUserData
	if err := ctl.request.DoValidate(c, &data); err != nil {
		return ctl.request.BadRequest(c, err)
	}

	ocservUser, err := ctl.ocservUserRepo.GetByUID(c.Request().Context(), userID)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}

	if data.Group != nil {
		ocservUser.Group = *data.Group
	}
	if data.Password != nil {
		ocservUser.Password = *data.Password
	}
	if data.Description != nil {
		ocservUser.Description = *data.Description
	}
	if data.TrafficSize != nil {
		ocservUser.TrafficSize = *data.TrafficSize
	}
	if data.TrafficType != nil && slices.Contains([]string{"Free", "MonthlyTransmit", "MonthlyReceive", "TotallyTransmit", "TotallyReceive"}, *data.TrafficType) {
		ocservUser.TrafficType = *data.TrafficType
	}
	if data.Config != nil {
		log.Println("\n\n")
		log.Println(data.Config)
		ocservUser.Config = data.Config
	}

	if data.ExpireAt != nil {
		expire, err := time.Parse("2006-01-02", *data.ExpireAt)
		if err != nil {
			expire, _ = time.Parse(
				"2006-01-02",
				time.Now().AddDate(0, 0, 30).Format("2006-01-02"),
			)
			ocservUser.ExpireAt = &expire
		}
	}

	updatedOcservUser, err := ctl.ocservUserRepo.Update(c.Request().Context(), ocservUser)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, updatedOcservUser)
}

// DeleteOcservUser 	     Ocserv User delete
//
// @Summary      Ocserv User delete
// @Description  Ocserv User delete
// @Tags         Ocserv(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      204  {object} nil
// @Router       /ocserv/users/{uid} [delete]
func (ctl *Controller) DeleteOcservUser(c echo.Context) error {
	userID := c.Param("uid")
	if userID == "" {
		return ctl.request.BadRequest(c, errors.New("user id is required"))
	}

	err := ctl.ocservUserRepo.Delete(c.Request().Context(), userID)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

// LockOcservUser 	     Ocserv User locking
//
// @Summary      Ocserv User locking
// @Description  Ocserv User locking
// @Tags         Ocserv(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object} nil
// @Router       /ocserv/users/{uid}/lock [post]
func (ctl *Controller) LockOcservUser(c echo.Context) error {
	userID := c.Param("uid")
	if userID == "" {
		return ctl.request.BadRequest(c, errors.New("user id is required"))
	}

	err := ctl.ocservUserRepo.Lock(c.Request().Context(), userID)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// UnLockOcservUser 	     Ocserv User unlocking
//
// @Summary      Ocserv User unlocking
// @Description  Ocserv User unlocking
// @Tags         Ocserv(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object} nil
// @Router       /ocserv/users/{uid}/unlock [post]
func (ctl *Controller) UnLockOcservUser(c echo.Context) error {
	userID := c.Param("uid")
	if userID == "" {
		return ctl.request.BadRequest(c, errors.New("user id is required"))
	}

	err := ctl.ocservUserRepo.UnLock(c.Request().Context(), userID)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// DisconnectOcservUser 	     Ocserv User disconnecting
//
// @Summary      Disconnect Ocserv User
// @Description  Disconnect Ocserv User
// @Tags         Ocserv(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 username path string true "Ocserv User username"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object} nil
// @Router       /ocserv/users/{username}/disconnect [post]
func (ctl *Controller) DisconnectOcservUser(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return ctl.request.BadRequest(c, errors.New("user id is required"))
	}
	_, err := ctl.ocRepo.Disconnect(c.Request().Context(), username)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, nil)
}

// StatisticsOcservUser 	     Ocserv User Statistics
//
// @Summary      Ocserv User Statistics
// @Description  Ocserv User Statistics
// @Tags         Ocserv(Users)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 uid path string true "Ocserv User UID"
// @Param 		 date_start query string false "date_start"
// @Param 		 date_end query string false "date_end"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200  {object} []models.DailyTraffic
// @Router       /ocserv/users/{uid}/statistics [get]
func (ctl *Controller) StatisticsOcservUser(c echo.Context) error {
	userID := c.Param("uid")
	if userID == "" {
		return ctl.request.BadRequest(c, errors.New("user id is required"))
	}

	var data StatisticsData
	if err := c.Bind(&data); err != nil {
		return ctl.request.BadRequest(c, err)
	}

	var startDate, endDate *time.Time

	if data.DateStart != "" {
		t, err := time.Parse("2006-01-02", data.DateStart)
		if err != nil {
			return ctl.request.BadRequest(c, fmt.Errorf("invalid date_start: %w", err))
		}
		startDate = &t
	}

	if data.DateEnd != "" {
		t, err := time.Parse("2006-01-02", data.DateEnd)
		if err != nil {
			return ctl.request.BadRequest(c, fmt.Errorf("invalid date_end: %w", err))
		}
		endDate = &t
	}

	stats, err := ctl.ocservUserRepo.UserStatistics(c.Request().Context(), userID, startDate, endDate)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, stats)
}

// Statistics 	 Ocserv Users Statistics
//
// @Summary      Ocserv Users Statistics
// @Description  Ocserv Users Statistics
// @Tags         Ocserv(Statistics)
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer TOKEN"
// @Param 		 date_start query string true "date_start"
// @Param 		 date_end query string true "date_end"
// @Failure      400 {object} request.ErrorResponse
// @Failure      401 {object} middlewares.Unauthorized
// @Success      200 {object} []models.DailyTraffic
// @Router       /ocserv/users/statistics [get]
func (ctl *Controller) Statistics(c echo.Context) error {
	var data StatisticsData
	if err := c.Bind(&data); err != nil {
		return ctl.request.BadRequest(c, err)
	}

	if data.DateStart == "" || data.DateEnd == "" {
		return ctl.request.BadRequest(c, errors.New("statistics date start and end are required"))
	}

	var startDate, endDate *time.Time

	t, err := time.Parse("2006-01-02", data.DateStart)
	if err != nil {
		return ctl.request.BadRequest(c, fmt.Errorf("invalid date_start: %w", err))
	}
	startDate = &t

	t, err = time.Parse("2006-01-02", data.DateEnd)
	if err != nil {
		return ctl.request.BadRequest(c, fmt.Errorf("invalid date_end: %w", err))
	}
	endDate = &t

	stats, err := ctl.ocservUserRepo.Statistics(c.Request().Context(), startDate, endDate)
	if err != nil {
		return ctl.request.BadRequest(c, err)
	}
	return c.JSON(http.StatusOK, stats)
}
