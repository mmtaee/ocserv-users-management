package home

import "ocserv-bakend/internal/models"

type StatsSections struct {
	GeneralInfo  string `json:"general_info" validate:"required"`
	CurrentStats string `json:"current_stats" validate:"required"`
}

type HomeResponse struct {
	Status     StatsSections          `json:"status" validate:"required"`
	Stats      *[]models.DailyTraffic `json:"stats" validate:"omitempty"`
	OnlineUser *[]string              `json:"online_user" validate:"omitempty"`
}
