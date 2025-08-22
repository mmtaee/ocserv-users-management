package system

import (
	"api/internal/models"
	"api/mocks"
	"api/pkg/crypto"
	"api/pkg/request"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
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

	ctl := &Controller{
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

// -------------------- TEST: SystemSetup --------------------
func TestSystemSetupSuccess(t *testing.T) {
	ctrl, mockRequest, mockSystemRepo, _, mockUserRepo, mockCryptoRepo := newControllerWithMocks()

	body := `{
		"username": "admin",
		"password": "test_password",
		"google_captcha_site_key": "",
		"google_captcha_secret_key": ""
	}`

	c, rec := setupEcho(http.MethodPost, "/system/setup", body)

	mockSystemRepo.On("System", mock.Anything).Return(nil, errors.New("not configured"))
	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).(*SetupSystem)
		data.Username = "admin"
		data.Password = "test_password"
		data.GoogleCaptchaSiteKey = ""
		data.GoogleCaptchaSecretKey = ""
	})
	mockCryptoRepo.On("CreatePassword", "test_password").Return(crypto.CustomPassword{
		Hash: "hashedPass",
		Salt: "saltValue",
	})

	user := &models.User{
		ID:       1,
		UID:      "uid-1234",
		Username: "admin",
		Password: "hashedPass",
		IsAdmin:  true,
	}

	system := &models.System{
		GoogleCaptchaSiteKey:   "",
		GoogleCaptchaSecretKey: "",
	}

	mockSystemRepo.On("SystemSetup", mock.Anything, mock.Anything, mock.Anything).Return(user, system, nil)
	mockUserRepo.On("CreateToken", mock.Anything, user.ID, user.UID, true, user.IsAdmin).Return("mock-token", nil)

	err := ctl.SetupSystem(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockSystemRepo.AssertExpectations(t)
	mockCryptoRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)

	var response SetupSystemResponse
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "admin", response.User.Username)
	assert.Equal(t, true, response.User.IsAdmin)
	assert.Equal(t, "", response.System.GoogleCaptchaSiteKey)
	assert.Equal(t, "", response.System.GoogleCaptchaSecretKey)
	assert.Equal(t, "mock-token", response.Token)
}

// -------------------- TEST: SystemInit --------------------
func TestSystemInitSuccess(t *testing.T) {
	ctrl, _, mockSystemRepo, _, _, _ := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/system/init", "")

	expected := &models.System{GoogleCaptchaSiteKey: "abc123"}
	mockSystemRepo.On("System", mock.Anything).Return(expected, nil)

	err := ctl.SystemInit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemInitNotFound(t *testing.T) {
	ctrl, _, mockSystemRepo, _, _, _ := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/system/init", "")

	mockSystemRepo.On("System", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := ctl.SystemInit(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockSystemRepo.AssertExpectations(t)
}

func TestSystemConfig(t *testing.T) {
	ctrl, _, mockSystemRepo, _, _, _ := newControllerWithMocks()

	c, rec := setupEcho(http.MethodGet, "/system", "")
	expected := &models.System{GoogleCaptchaSiteKey: "abc123", GoogleCaptchaSecretKey: "abc123"}
	mockSystemRepo.On("System", mock.Anything).Return(expected, nil)
	err := ctl.System(c)
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

	err := ctl.SystemUpdate(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRequest.AssertExpectations(t)
	mockSystemRepo.AssertExpectations(t)
}

// -------------------- TEST: Login --------------------
func TestSystemLogin(t *testing.T) {
	ctrl, mockRequest, mockSystemRepo, mockCaptcha, mockUserRepo, mockCryptoRepo := newControllerWithMocks()

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

	mockUser := &models.User{ID: 1, UID: "uid-123", Username: "testuser", IsAdmin: false, Password: "testpass", Salt: "testsalt"}
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(mockUser, nil)

	mockCryptoRepo.On("CheckPassword", "testpass", mockUser.Password, mockUser.Salt).Return(true)

	mockUserRepo.On("CreateToken", mock.Anything, uint(1), "uid-123", true, false).Return("mock-token", nil)

	err := ctl.Login(c)

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

	err := ctl.CreateUser(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "uid-123", resp["uid"])
	assert.Equal(t, "testuser", resp["username"])
	assert.Equal(t, false, resp["is_admin"])

	assert.Contains(t, resp, "CreatedAt")
	assert.Contains(t, resp, "UpdatedAt")

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

	err := ctl.Users(c)

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

func TestChangeUserPasswordByAdmin(t *testing.T) {
	ctrl, mockRequest, _, _, mockUserRepo, mockCryptoRepo := newControllerWithMocks()
	passwordInput := `{"password": "testpass"}`

	c, rec := setupEcho(http.MethodPost, "/system/users/uid-123/password", passwordInput)
	c.SetPath("/system/users/:id/password")
	c.SetParamNames("uid")
	c.SetParamValues("uid-123")
	c.Set("isAdmin", true)

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).(*ChangeUserPassword)
		data.Password = "testpass"
	})

	mockCryptoRepo.On("CreatePassword", "testpass").Return(crypto.CustomPassword{
		Hash: "hashedPass",
		Salt: "saltValue",
	})

	mockUserRepo.On("ChangePassword", mock.Anything, "uid-123", "hashedPass", "saltValue").Return(nil)

	err := ctl.ChangeUserPasswordByAdmin(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockRequest.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockCryptoRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	ctrl, _, _, _, mockUserRepo, _ := newControllerWithMocks()
	c, rec := setupEcho(http.MethodDelete, "/system/users/uid-123", "")
	c.SetPath("/system/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("uid-123")
	c.Set("isAdmin", true)

	mockUserRepo.On("DeleteUser", mock.Anything, mock.Anything).Return(nil)
	err := ctl.DeleteUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUserRepo.AssertExpectations(t)
}

func TestChangePasswordBySelf(t *testing.T) {
	ctrl, mockRequest, _, _, mockUserRepo, mockCryptoRepo := newControllerWithMocks()
	passwordInput := `{"password": "testpass"}`

	c, rec := setupEcho(http.MethodPost, "/system/users/password", passwordInput)
	c.SetPath("/system/users/password")
	c.Set("isAdmin", true)
	c.Set("userUID", "uid-123")

	mockRequest.On("DoValidate", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		data := args.Get(1).(*ChangeUserPassword)
		data.Password = "testpass"
	})

	mockCryptoRepo.On("CreatePassword", "testpass").Return(crypto.CustomPassword{
		Hash: "hashedPass",
		Salt: "saltValue",
	})

	mockUserRepo.On("ChangePassword", mock.Anything, "uid-123", "hashedPass", "saltValue").Return(nil)

	err := ctl.ChangePasswordBySelf(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	mockRequest.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockCryptoRepo.AssertExpectations(t)
}

func TestUserProfile(t *testing.T) {
	uid := "uid-1234"
	ctrl, _, _, _, mockUserRepo, _ := newControllerWithMocks()
	c, rec := setupEcho(http.MethodGet, "/system/users/profile", "")
	c.Set("userUID", uid)

	mockUserRepo.On("GetByUID", mock.Anything, uid).Return(&models.User{
		UID:      uid,
		Username: "testuser",
		IsAdmin:  true,
	}, nil)

	err := ctl.Profile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUserRepo.AssertExpectations(t)

	var response models.User
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, uid, response.UID)
	assert.Equal(t, true, response.IsAdmin)
	assert.Equal(t, "testuser", response.Username)
	mockUserRepo.AssertExpectations(t)
}
