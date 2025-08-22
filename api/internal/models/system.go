package models

import (
	"api/pkg/database"
	"errors"
	"gorm.io/gorm"
)

type System struct {
	ID                     uint   `json:"_" gorm:"primaryKey"`
	GoogleCaptchaSecretKey string `json:"google_captcha_secret" gorm:"type:text"`
	GoogleCaptchaSiteKey   string `json:"google_captcha_site_key" gorm:"type:text"`
}

func (s *System) BeforeCreate(tx *gorm.DB) error {
	ch := make(chan error, 1)
	go func() {
		var system System
		db := database.Get()
		err := db.Table("panels").First(&system).Error
		if err != nil && system.ID == 0 {
			ch <- nil
		}
		ch <- errors.New("system configs already exist")
	}()
	return <-ch
}
