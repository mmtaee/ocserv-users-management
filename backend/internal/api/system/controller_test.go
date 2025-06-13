package system

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/mocks"
	"ocserv-bakend/pkg/crypto"
	"ocserv-bakend/pkg/request"
	"strings"
	"testing"
)

func newControllerWithMocks() (
	*Controller,
	*mocks.CustomRequestInterface,
	*mocks.SystemRepositoryInterface,
	*mocks.GoogleCaptchaInterface,
	*mocks.UserRepositoryInterface,
	*mocks.CustomPasswordInterface,
) {
	mockRequest := new(mocks.CustomRequestInterface)
	mockSystemRepo := new(mocks.SystemRepositoryInterface)
	mockCaptcha := new(mocks.GoogleCaptchaInterface)
	mockUserRepo := new(mocks.UserRepositoryInterface)
	mockCryptoRepo := new(mocks.CustomPasswordInterface)

	ctrl := &Controller{
		request:         mockRequest,
		systemRepo:      mockSystemRepo,
		captchaVerifier: mockCaptcha,
		userRepo:        mockUserRepo,
		cryptoRepo:      mockCryptoRepo,
	}

	return ctrl, mockRequest, mockSystemRepo, mockCaptcha, mockUserRepo, mockCryptoRepo
}

func strPtr(s string) *string {
	return &s
}

func setupEcho(method, path string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// -------------------- TEST: SystemInit --------------------
func TestSystemInitSuccess(t *testing.T) {
	ctrl, _, mockSystemRepo, _, _, _ := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/system/init", "")

	expected := &models.System{GoogleCaptchaSiteKey: "abc123"}
	mockSystemRepo.On("System", mock.Anything).Return(expected, nil)

	err := ctrl.SystemInit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemInitNotFound(t *testing.T) {
	ctrl, _, mockSystemRepo, _, _, _ := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/system/init", "")

	mockSystemRepo.On("System", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := ctrl.SystemInit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemConfig(t *testing.T) {
	ctrl, _, mockSystemRepo, _, _, _ := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/system", "")
	expected := &models.System{GoogleCaptchaSiteKey: "abc123", GoogleCaptchaSecretKey: "abc123"}
	mockSystemRepo.On("System", mock.Anything).Return(expected, nil)
	err := ctrl.System(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockSystemRepo.AssertExpectations(t)
}

// -------------------- TEST: SystemUpdate --------------------
func TestSystemUpdateSuccess(t *testing.T) {
	ctrl, mockRequest, mockSystemRepo, _, _, _ := newControllerWithMocks()

	body := `{"google_captcha_site_key":"key123","google_captcha_secret_key":"secret456"}`
	c, rec := setupEcho(http.MethodPatch, "/system", body)

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

// -------------------- TEST: Login --------------------
func TestSystemLogin(t *testing.T) {
	ctrl, mockRequest, mockSystemRepo, mockCaptcha, mockUserRepo, _ := newControllerWithMocks()

	loginInput := `{"username":"testuser", "password":"testpass", "token":"dummy-token"}`
	c, w := setupEcho(http.MethodPost, "/system/users/login", loginInput)

	loginData := &LoginData{
		Username: "testuser",
		Password: "testpass",
		Token:    "dummy-token",
	}

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*LoginData)
		*arg = *loginData
	})

	config := &models.System{GoogleCaptchaSecretKey: "captcha-key", GoogleCaptchaSiteKey: "abc123"}
	mockSystemRepo.On("System", mock.Anything).Return(config, nil)

	mockCaptcha.On("SetSecretKey", "captcha-key").Return(mockCaptcha)
	mockCaptcha.On("Verify", "dummy-token").Return(mockCaptcha)
	mockCaptcha.On("IsValid").Return(true)

	mockUser := &models.User{ID: 1, UID: "uid-123", Username: "testuser", IsAdmin: false}
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(mockUser, nil)
	mockUserRepo.On("CreateToken", mock.Anything, uint(1), "uid-123", true, false).Return("mock-token", nil)

	err := ctrl.Login(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRequest.AssertExpectations(t)
	mockSystemRepo.AssertExpectations(t)
	mockCaptcha.AssertExpectations(t)
	mockSystemRepo.AssertExpectations(t)
}

// -------------------- TEST: user --------------------
func TestCreateUserSuccess(t *testing.T) {
	ctrl, mockRequest, _, _, mockUserRepo, mockCryptoRepo := newControllerWithMocks()

	userInput := `{"username":"testuser", "password":"testpass", "admin":false}`
	c, rec := setupEcho(http.MethodPost, "/system/users", userInput)

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).(*CreateUserData)
		data.Username = "testuser"
		data.Password = "testpass"
		data.Admin = false
	})

	mockCryptoRepo.On("CreatePassword", "testpass").Return(crypto.CustomPassword{
		Hash: "hashedPass",
		Salt: "saltValue",
	})

	mockUser := &models.User{
		ID:       1,
		UID:      "uid-123",
		Username: "testuser",
		Password: "hashedPass",
		IsAdmin:  false,
	}

	mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).Return(mockUser, nil)

	err := ctrl.CreateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), resp["id"])
	assert.Equal(t, "uid-123", resp["uid"])
	assert.Equal(t, "testuser", resp["username"])
	assert.Equal(t, false, resp["is_admin"])

	assert.Contains(t, resp, "CreatedAt")
	assert.Contains(t, resp, "UpdatedAt")
	assert.Contains(t, resp, "last_login")

	mockRequest.AssertExpectations(t)
	mockCryptoRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestUsers(t *testing.T) {
	ctrl, mockRequest, _, _, mockUserRepo, _ := newControllerWithMocks()
	c, rec := setupEcho(http.MethodGet, "/system/users", "")

	pagination := &request.Pagination{Page: 1, PageSize: 10}
	mockRequest.On("Pagination", mock.AnythingOfType("*echo.context")).Return(pagination)

	mockUserRepo.
		On("Users", mock.Anything, pagination).
		Return(&[]models.User{}, int64(0), nil)

	err := ctrl.Users(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp UsersResponse
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), resp.Meta.TotalRecords)
	assert.Empty(t, resp.Result)
	
	mockRequest.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
