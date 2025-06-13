package system

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/mocks"
	"ocserv-bakend/pkg/request"
	"strings"
	"testing"
)

var (
	mockRequest    = new(mocks.CustomRequestInterface)
	mockSystemRepo = new(mocks.SystemRepositoryInterface)

	ctrl = &Controller{
		request:    mockRequest,
		systemRepo: mockSystemRepo,
	}
)

func strPtr(s string) *string {
	return &s
}

// -------------------- TEST: SystemInit --------------------
func TestSystemInitSuccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/system/Init", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expected := &models.System{GoogleCaptchaSiteKey: "abc123"}
	mockSystemRepo.On("System", mock.Anything).Return(expected, nil)

	err := ctrl.SystemInit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemInitNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/system/Init", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockSystemRepo.On("System", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := ctrl.SystemInit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockSystemRepo.AssertExpectations(t)
}

// -------------------- TEST: SystemUpdate --------------------
func TestSystemUpdateSuccess(t *testing.T) {
	e := echo.New()
	body := `{"googleCaptchaSiteKey":"key123","googleCaptchaSecretKey":"secret456"}`
	req := httptest.NewRequest(http.MethodPatch, "/system", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*PatchSystemUpdateData)
		arg.GoogleCaptchaSiteKey = strPtr("key123")
		arg.GoogleCaptchaSecretKey = strPtr("secret456")
	})

	expected := &models.System{
		GoogleCaptchaSiteKey:   "key123",
		GoogleCaptchaSecretKey: "secret456",
	}
	mockSystemRepo.On("SystemUpdate", mock.Anything, mock.Anything).Return(expected, nil)

	err := ctrl.SystemUpdate(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRequest.AssertExpectations(t)
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemUpdateValidationFail(t *testing.T) {
	body := `{"googleCaptchaSiteKey":"", "googleCaptchaSecretKey":""}`

	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/system", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	validationErr := errors.New("validation error")

	// ✅ Setup mock expectations
	mockRequest.On("DoValidate", mock.Anything, mock.Anything).
		Return(validationErr).Once()

	mockRequest.On("BadRequest", mock.AnythingOfType("*echo.context"), validationErr).
		Return(func(c echo.Context, err error) error {
			return c.JSON(http.StatusBadRequest, request.ErrorResponse{
				Error:   []string{err.Error()},
				Message: []string{err.Error()},
			})
		}).Once()
	err := ctrl.SystemUpdate(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	//mockRequest.AssertExpectations(t)

	var resp request.ErrorResponse
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	assert.Contains(t, resp.Message, "validation error")
}
