package log

import (
	auditLogs "ocserv-bakend/pkg/audit_log"
	"ocserv-bakend/pkg/request"
)

type UsersLogsResponse struct {
	Meta   request.Meta          `json:"meta" validate:"required"`
	Result *[]auditLogs.AuditLog `json:"result" validate:"omitempty"`
}
