package ocserv_group

import (
	"encoding/json"
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

func TestOcservGroupList(t *testing.T) {
	ctrl, mockRequest, mockOcservUserRepo := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/ocserv/groups", "")

	pagination := &request.Pagination{Page: 1, PageSize: 10}
	mockRequest.On("Pagination", mock.AnythingOfType("*echo.context")).Return(pagination)

	mockOcservUserRepo.
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
	mockOcservUserRepo.AssertExpectations(t)

}
