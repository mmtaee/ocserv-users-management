package system

type GetSystemInitResponse struct {
	GoogleCaptchaSiteKey string `json:"google_captcha_site_key" validate:"omitempty"`
}

type GetSystemResponse struct {
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" validate:"omitempty"`
	GoogleCaptchaSecretKey string `json:"google_captcha_secret_key" validate:"omitempty"`
}

type PatchSystemUpdateData struct {
	GoogleCaptchaSiteKey   *string `json:"google_captcha_site_key" validate:"omitempty"`
	GoogleCaptchaSecretKey *string `json:"google_captcha_secret_key" validate:"omitempty"`
}
