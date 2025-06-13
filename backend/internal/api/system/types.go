package system

import "ocserv-bakend/internal/models"

type GetSystemInitResponse struct {
	GoogleCaptchaSiteKey string `json:"google_captcha_site_key" validate:"omitempty"`
}

type GetSystemResponse struct {
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" validate:"omitempty"`
	GoogleCaptchaSecretKey string `json:"google_captcha_secret_key" validate:"omitempty"`
}

type PatchSystemUpdateData struct {
	GoogleCaptchaSiteKey   *string `json:"google_captcha_site_key" validate:"required"`
	GoogleCaptchaSecretKey *string `json:"google_captcha_secret_key" validate:"required"`
}

type LoginData struct {
	Username   string `json:"username" validate:"required,min=2,max=16" example:"john_doe" `
	Password   string `json:"password" validate:"required,min=2,max=16" example:"doe123456"`
	RememberMe bool   `json:"remember_me" desc:"remember for a month"`
	Token      string `json:"token" desc:"captcha v2 token"`
}

type UserLoginResponse struct {
	User  *models.User `json:"user" validate:"required"`
	Token string       `json:"token" validate:"required"`
}

type CreateUserData struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Admin    bool   `json:"admin"`
}
