package ocserv_user

import (
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/request"
)

type CreateOcservUserData struct {
	Group       string `json:"group" validate:"required"`
	Username    string `json:"username" validate:"required,min=2,max=32"`
	Password    string `json:"password" validate:"required,min=2,max=32"`
	ExpireAt    string `json:"expire_at" validate:"omitempty" example:"2025-12-31"`
	TrafficType string `json:"traffic_type" validate:"required,oneof=Free MonthlyTransmit MonthlyReceive TotallyTransmit TotallyReceive" example:"MonthlyTransmit"`
	TrafficSize int    `json:"traffic_size" validate:"omitempty,gte=0" example:"10737418240"` // 10 GiB
	Description string `json:"description,omitempty" validate:"omitempty,max=1024" example:"User for testing VPN access"`
}

type UpdateOcservUserData struct {
	Group       *string `json:"group" example:"default"`
	Password    *string `json:"password" validate:"min=6" example:"strongpassword123"`
	ExpireAt    *string `json:"expire_at"  example:"2025-12-31"`
	TrafficType *string `json:"traffic_type" validate:"oneof=Free MonthlyTransmit MonthlyReceive TotallyTransmit TotallyReceive" example:"MonthlyTransmit"`
	TrafficSize *int    `json:"traffic_size" validate:"gt=0" example:"10737418240"` // 10 GiB
	Description *string `json:"description,omitempty" validate:"omitempty,max=1024" example:"User for testing VPN access"`
}

type OcservUsersResponse struct {
	Meta   request.Meta         `json:"meta" validate:"required"`
	Result *[]models.OcservUser `json:"result" validate:"omitempty"`
}
