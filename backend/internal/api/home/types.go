package home

import "ocserv-bakend/internal/models"

type StatsSections struct {
	GeneralInfo  string `json:"general_info" validate:"required"`
	CurrentStats string `json:"current_stats" validate:"required"`
}

type HomeResponse struct {
	Status     StatsSections               `json:"status" validate:"required"`
	Stats      *[]models.DailyTraffic      `json:"stats" validate:"omitempty"`
	OnlineUser *[]models.OnlineUserSession `json:"online_users_session" validate:"omitempty"`
	IPBans     *[]models.IPBan             `json:"ipbans" validate:"omitempty"`
	//IRoutes    *[]models.Iroute       `json:"iroutes" validate:"omitempty"` // has bug on version 1.2.4
}
