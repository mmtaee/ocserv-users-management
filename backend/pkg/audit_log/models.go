package audit_log

import "time"

const (
	EventCreate = "create"
	EventUpdate = "update"
	EventDelete = "delete"

	EventUserModel        = "user"
	EventOcservGroupModel = "ocserv_group"
)

type AuditLogAction struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
}

type AuditLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserUID   string    `json:"user_uid" gorm:"type:varchar(26);index"` // ULID or user identifier
	Model     string    `json:"model" gorm:"type:varchar(16);index"`    // e.g., "OcservUser"
	ModelID   string    `json:"model_id" gorm:"type:varchar(26);index"` // e.g., UID or primary key
	Action    string    `json:"action" gorm:"type:varchar(128)"`        // action and reason json string
	Changes   string    `json:"changes" gorm:"type:json"`               // Full JSON or JSON diff
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
