package ocserv_group

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
)

func newControllerWithMocks() (
	*Controller,
	*mocks.CustomRequestInterface,
	*mocks.OcservGroupRepositoryInterface,
) {
	mockRequest := new(mocks.CustomRequestInterface)
	ocservGroupRepo := new(mocks.OcservGroupRepositoryInterface)

	ctl := &Controller{
		request:         mockRequest,
		ocservGroupRepo: ocservGroupRepo,
	}

	return ctrl, mockRequest, ocservGroupRepo
}

func setupEcho(method, path string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func intPtr(i int) *int {
	return &i
}

func TestOcservGroupList(t *testing.T) {
	ctrl, mockRequest, ocservGroupRepo := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/ocserv/groups", "")

	pagination := &request.Pagination{Page: 1, PageSize: 10}
	mockRequest.On("Pagination", mock.AnythingOfType("*echo.context")).Return(pagination)

	ocservGroupRepo.
		On("Groups", mock.Anything, pagination).
		Return(&[]models.OcservGroup{}, int64(0), nil)

	err := ctl.OcservGroups(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp OcservGroupsResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), resp.Meta.TotalRecords)
	assert.Empty(t, resp.Result)

	mockRequest.AssertExpectations(t)
	ocservGroupRepo.AssertExpectations(t)
}

func TestOcservGroupLookup(t *testing.T) {
	ctrl, _, ocservGroupRepo := newControllerWithMocks()
	c, rec := setupEcho(http.MethodGet, "/ocserv/groups/lookup", "")

	ocservGroupRepo.On("GroupsLookup", mock.Anything).Return([]string{}, nil)
	err := ctl.OcservGroupsLookup(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	var resp []string
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, []string{"defaults"}, resp)

	ocservGroupRepo.AssertExpectations(t)
}

func TestCreateOcservGroupSuccess(t *testing.T) {
	ctrl, mockRequest, ocservGroupRepo := newControllerWithMocks()

	mockInput := CreateOcservGroupData{
		Name: "test-group",
		Config: &models.OcservGroupConfig{
			MTU:            intPtr(1300),
			MaxSameClients: intPtr(3),
		},
	}

	body, _ := json.Marshal(mockInput)

	c, rec := setupEcho(http.MethodGet, "/ocserv/groups", string(body))

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).(*CreateOcservGroupData)
		*data = mockInput
	})

	expectedGroup := models.OcservGroup{
		ID:   1,
		UID:  "uid-123",
		Name: "test-group",
		Config: &models.OcservGroupConfig{
			MTU:            intPtr(1300),
			MaxSameClients: intPtr(3),
		},
	}

	ocservGroupRepo.On("Create", mock.Anything, mock.Anything).
		Return(&expectedGroup, nil)

	err := ctl.CreateOcservGroup(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "test-group", resp["name"])
	assert.Equal(t, "uid-123", resp["uid"])

	mockRequest.AssertExpectations(t)
	ocservGroupRepo.AssertExpectations(t)
}

func TestCreateOcservGroupFailed(t *testing.T) {
	ctrl, mockRequest, _ := newControllerWithMocks()
	c, rec := setupEcho(http.MethodGet, "/ocserv/groups", "")

	expectedErr := errors.New("validation error")
	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(expectedErr)

	mockRequest.On("BadRequest", mock.Anything, expectedErr).
		Return(c.JSON(http.StatusBadRequest, expectedErr))

	err := ctl.CreateOcservGroup(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRequest.AssertExpectations(t)
}

func TestUpdateOcservGroupSuccess(t *testing.T) {
	ctrl, mockRequest, ocservGroupRepo := newControllerWithMocks()

	mockInput := UpdateOcservGroupData{
		Config: &models.OcservGroupConfig{
			MTU:            intPtr(1300),
			MaxSameClients: intPtr(3),
		},
	}

	body, _ := json.Marshal(mockInput)

	c, rec := setupEcho(http.MethodPatch, "/ocserv/groups/uid-1234", string(body))

	c.SetPath("/ocserv/groups/:uid")
	c.SetParamNames("uid")
	c.SetParamValues("uid-123")

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).(*UpdateOcservGroupData)
		*data = mockInput
	})

	ocservGroup := models.OcservGroup{
		ID:   1,
		UID:  "uid-1234",
		Name: "test-group",
		Config: &models.OcservGroupConfig{
			MTU:            intPtr(1300),
			MaxSameClients: intPtr(3),
		},
	}

	ocservGroupRepo.On("GetByUID", mock.Anything, "uid-123").Return(&ocservGroup, nil)

	updatedGroup := ocservGroup
	updatedGroup.Config.MTU = intPtr(1400)
	updatedGroup.Config.MaxSameClients = intPtr(3)
	ocservGroupRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.OcservGroup")).Return(&updatedGroup, nil)

	err := ctl.UpdateOcservGroup(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.OcservGroup
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1400, *resp.Config.MTU)
	assert.Equal(t, 3, *resp.Config.MaxSameClients)

	mockRequest.AssertExpectations(t)
	ocservGroupRepo.AssertExpectations(t)
}

func TestUpdateOcservGroupFailed(t *testing.T) {
	ctrl, mockRequest, _ := newControllerWithMocks()

	c, rec := setupEcho(http.MethodPatch, "/ocserv/groups/uid-1234", "")

	c.SetPath("/ocserv/groups/:uid")
	c.SetParamNames("uid")
	c.SetParamValues("uid-123")

	expectedErr := errors.New("validation error")
	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(expectedErr)

	mockRequest.On("BadRequest", mock.Anything, expectedErr).
		Return(c.JSON(http.StatusBadRequest, expectedErr))

	err := ctl.UpdateOcservGroup(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRequest.AssertExpectations(t)
}

func TestDeleteOcservGroupSuccess(t *testing.T) {
	ctrl, _, ocservGroupRepo := newControllerWithMocks()

	c, rec := setupEcho(http.MethodDelete, "/ocserv/groups/uid-1234", "")

	c.SetPath("/ocserv/groups/:uid")
	c.SetParamNames("uid")
	c.SetParamValues("uid-123")

	ocservGroupRepo.On("Delete", mock.Anything, "uid-123").Return(nil)

	err := ctl.DeleteOcservGroup(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	ocservGroupRepo.AssertExpectations(t)
}
