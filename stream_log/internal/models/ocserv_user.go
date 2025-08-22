package models

import (
	"time"
)

const (
	Free            = "Free"
	MonthlyTransmit = "MonthlyTransmit"
	MonthlyReceive  = "MonthlyReceive"
	TotallyTransmit = "TotallyTransmit"
	TotallyReceive  = "TotallyReceive"
)

type OcservUser struct {
	ID            uint       `json:"-" gorm:"primaryKey;autoIncrement" `
	UID           string     `json:"uid" gorm:"type:varchar(26);not null;unique" validate:"required"`
	Username      string     `json:"username" gorm:"type:varchar(16);not null;unique" validate:"required"`
	IsLocked      bool       `json:"is_locked" gorm:"default(false)" validate:"required"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime" validate:"required"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime" validate:"omitempty"`
	ExpireAt      *time.Time `json:"expire_at" gorm:"type:date" validate:"omitempty"`
	DeactivatedAt *time.Time `json:"deactivated_at" gorm:"type:date" validate:"omitempty"`
	TrafficType   string     `json:"traffic_type" gorm:"type:varchar(32);not null;default:1" enums:"Free,MonthlyTransmit,MonthlyReceive,TotallyTransmit,TotallyReceive" validate:"required"`
	TrafficSize   int        `json:"traffic_size" gorm:"not null" validate:"required"` // in GiB  >> x * 1024 ** 3
	Rx            int        `json:"rx" gorm:"not null;default:0" validate:"required"` // Receive in bytes
	Tx            int        `json:"tx" gorm:"not null;default:0" validate:"required"` // Transmit in bytes
}

type OcservUserTrafficStatistics struct {
	ID        uint      `json:"-" gorm:"primaryKey;autoIncrement"`
	OcUserID  uint      `json:"-" gorm:"index;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	Rx        int       `json:"rx" gorm:"default:0"` // in bytes
	Tx        int       `json:"tx" gorm:"default:0"` // in bytes
}
