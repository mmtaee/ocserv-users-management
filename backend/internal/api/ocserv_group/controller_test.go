package ocserv_group

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/mocks"
	"ocserv-bakend/pkg/request"
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

	ctrl := &Controller{
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

func strPtr(s string) *string {
	return &s
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

	err := ctrl.OcservGroups(c)

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

	err := ctrl.CreateOcservGroup(c)

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

	err := ctrl.CreateOcservGroup(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockRequest.AssertExpectations(t)
}
