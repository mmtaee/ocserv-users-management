package log

import (
	auditLogs "api/pkg/audit_log"
	"api/pkg/request"
)

type UsersLogsResponse struct {
	Meta   request.Meta          `json:"meta" validate:"required"`
	Result *[]auditLogs.AuditLog `json:"result" validate:"omitempty"`
}
