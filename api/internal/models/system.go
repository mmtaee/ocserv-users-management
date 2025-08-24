package models

import (
	"errors"
	"github.com/mmtaee/ocserv-users-management/common/pkg/database"
	"gorm.io/gorm"
)

type System struct {
	ID                     uint   `json:"_" gorm:"primaryKey"`
	GoogleCaptchaSecretKey string `json:"google_captcha_secret" gorm:"type:text"`
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" gorm:"type:text"`
}

func (s *System) BeforeCreate(tx *gorm.DB) error {
	var system System
	db := database.GetConnection()
	err := db.Table("systems").First(&system).Error
	if err != nil && system.ID == 0 {
		return errors.New("system configs already exist")
	}
	return nil
}
