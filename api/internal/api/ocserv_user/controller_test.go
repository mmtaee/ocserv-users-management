package ocserv_user

import (
	"api/internal/models"
	"api/mocks"
	"api/pkg/request"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newControllerWithMocks() (
	*Controller,
	*mocks.CustomRequestInterface,
	*mocks.OcservUserRepositoryInterface,
) {
	mockRequest := new(mocks.CustomRequestInterface)
	mockOcservUserRepo := new(mocks.OcservUserRepositoryInterface)

	ctl := &Controller{
		request:        mockRequest,
		ocservUserRepo: mockOcservUserRepo,
	}

	return ctrl, mockRequest, mockOcservUserRepo
}

func setupEcho(method, path string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func mustParseTime(value string) time.Time {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}
	return t
}

func strPtr(s string) *string {
	return &s
}

func TestOcservUsers(t *testing.T) {
	ctrl, mockRequest, mockOcservUserRepo := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/ocserv/users", "")

	pagination := &request.Pagination{Page: 1, PageSize: 10}
	mockRequest.On("Pagination", mock.AnythingOfType("*echo.context")).Return(pagination)

	mockOcservUserRepo.
		On("Users", mock.Anything, pagination).
		Return(&[]models.OcservUser{}, int64(0), nil)

	err := ctl.OcservUsers(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp OcservUsersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), resp.Meta.TotalRecords)
	assert.Empty(t, resp.Result)

	mockRequest.AssertExpectations(t)
	mockOcservUserRepo.AssertExpectations(t)
}

func TestOcservUserCreateSuccess(t *testing.T) {
	ctrl, mockRequest, mockOcservUserRepo := newControllerWithMocks()

	expireAt := "2025-12-31"

	mockInput := CreateOcservUserData{
		Username:    "testuser",
		Password:    "testpass",
		Group:       "defaults",
		ExpireAt:    expireAt,
		TrafficSize: 0,
		TrafficType: models.Free,
	}

	body, _ := json.Marshal(mockInput)
	c, rec := setupEcho(http.MethodPost, "/ocserv/users", string(body))

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).(*CreateOcservUserData)
		*data = mockInput
	})

	expectedUser := models.OcservUser{
		ID:          1,
		Username:    mockInput.Username,
		Password:    mockInput.Password,
		Group:       mockInput.Group,
		ExpireAt:    mustParseTime(expireAt),
		TrafficSize: mockInput.TrafficSize,
		TrafficType: mockInput.TrafficType,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockOcservUserRepo.On("Create", mock.Anything, mock.Anything).
		Return(&expectedUser, nil)

	err := ctl.CreateOcservUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "defaults", resp["group"])
	assert.Equal(t, "testuser", resp["username"])
	assert.Equal(t, "testpass", resp["password"])
	assert.Contains(t, resp, "created_at")
	assert.Contains(t, resp, "updated_at")

	mockRequest.AssertExpectations(t)
	mockOcservUserRepo.AssertExpectations(t)
}

func TestOcservUserCreateFailed(t *testing.T) {
	ctrl, mockRequest, _ := newControllerWithMocks()
	c, rec := setupEcho(http.MethodPost, "/ocserv/users", "")

	expectedErr := errors.New("validation error")
	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(expectedErr)

	mockRequest.On("BadRequest", mock.Anything, expectedErr).
		Return(c.JSON(http.StatusBadRequest, expectedErr))

	err := ctl.CreateOcservUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRequest.AssertExpectations(t)
}

func TestOcservUserUpdateSuccess(t *testing.T) {
	ctrl, mockRequest, mockOcservUserRepo := newControllerWithMocks()

	mockInput := UpdateOcservUserData{
		Password: strPtr("updated-password"),
	}
	body, _ := json.Marshal(mockInput)
	c, rec := setupEcho(http.MethodPut, "/ocserv/users/uid-123", string(body))

	c.SetPath("/ocserv/users/:uid")
	c.SetParamNames("uid")
	c.SetParamValues("uid-123")

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).(*UpdateOcservUserData)
		*data = mockInput
	})

	ocservUser := models.OcservUser{
		ID:       1,
		UID:      "uid-123",
		Username: "testuser",
		Password: "testpass",
		Group:    "defaults",
	}

	mockOcservUserRepo.On("GetByUID", mock.Anything, "uid-123").Return(&ocservUser, nil)

	updatedUser := ocservUser
	updatedUser.Password = "updated-password"
	mockOcservUserRepo.On("Update", mock.Anything, &ocservUser).Return(&updatedUser, nil)
	mockOcservUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.OcservUser")).Return(&updatedUser, nil)

	err := ctl.UpdateOcservUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.OcservUser
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "updated-password", resp.Password)

	mockRequest.AssertExpectations(t)
	mockOcservUserRepo.AssertExpectations(t)
}

func TestDeleteOcservUserSuccess(t *testing.T) {
	ctrl, _, mockOcservUserRepo := newControllerWithMocks()

	c, rec := setupEcho(http.MethodDelete, "/ocserv/users/uid-123", "")

	c.SetPath("/ocserv/users/:uid")
	c.SetParamNames("uid")
	c.SetParamValues("uid-123")

	mockOcservUserRepo.On("Delete", mock.Anything, "uid-123").Return(nil)

	err := ctl.DeleteOcservUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockOcservUserRepo.AssertExpectations(t)
}

func TestDeleteOcservUserFailed(t *testing.T) {
	ctrl, mockRequest, _ := newControllerWithMocks()

	c, rec := setupEcho(http.MethodDelete, "/ocserv/users/", "")

	c.SetPath("/ocserv/users/:uid")
	c.SetParamNames("uid")
	c.SetParamValues("")

	expectedErr := errors.New("user id is required")

	mockRequest.On("BadRequest", mock.Anything, expectedErr).
		Return(c.JSON(http.StatusBadRequest, expectedErr))

	err := ctl.DeleteOcservUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRequest.AssertExpectations(t)
}

func TestLockOcservUser(t *testing.T) {
	ctrl, _, mockOcservUserRepo := newControllerWithMocks()
	c, rec := setupEcho(http.MethodPost, "/ocserv/users/uid-123/lock", "")
	c.SetPath("/ocserv/users/:uid/lock")
	c.SetParamNames("uid")
	c.SetParamValues("uid-123")

	mockOcservUserRepo.On("Lock", mock.Anything, "uid-123").Return(nil)
	err := ctl.LockOcservUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockOcservUserRepo.AssertExpectations(t)
}

func TestUnLockOcservUser(t *testing.T) {
	ctrl, _, mockOcservUserRepo := newControllerWithMocks()
	c, rec := setupEcho(http.MethodPost, "/ocserv/users/uid-123/unlock", "")
	c.SetPath("/ocserv/users/:uid/unlock")
	c.SetParamNames("uid")
	c.SetParamValues("uid-123")

	mockOcservUserRepo.On("UnLock", mock.Anything, "uid-123").Return(nil)
	err := ctl.UnLockOcservUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockOcservUserRepo.AssertExpectations(t)
}
