package repository

import (
	auditLogs "api/pkg/audit_log"
	"api/pkg/database"
	"api/pkg/request"
	"context"
	"gorm.io/gorm"
)

type LogsRepository struct {
	db *gorm.DB
}
type LogsRepositoryInterface interface {
	UsersLogs(ctx context.Context, pagination *request.Pagination, userUID string) (*[]auditLogs.AuditLog, int64, error)
	Logs(ctx context.Context, pagination *request.Pagination, userUID ...string) (*[]auditLogs.AuditLog, int64, error)
}

func NewLogsRepository() *LogsRepository {
	return &LogsRepository{
		db: database.Get(),
	}
}

func (l *LogsRepository) UsersLogs(ctx context.Context, pagination *request.Pagination, userUID string) (*[]auditLogs.AuditLog, int64, error) {
	var (
		totalRecords int64
		logs         []auditLogs.AuditLog
	)

	filter := "model = ? AND user_uid = ?"
	args := []interface{}{auditLogs.EventUserModel, userUID}

	if err := l.db.WithContext(ctx).
		Model(&auditLogs.AuditLog{}).
		Where(filter, args...).
		Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	txPaginator := paginator(ctx, l.db, pagination)
	err := txPaginator.Model(&logs).
		Where(filter, args...).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}
	return &logs, totalRecords, nil
}

func (l *LogsRepository) Logs(ctx context.Context, pagination *request.Pagination, userUID ...string) (*[]auditLogs.AuditLog, int64, error) {
	var (
		totalRecords int64
		logs         []auditLogs.AuditLog
		filters      string
		args         []interface{}
	)

	if len(userUID) > 0 && userUID[0] != "" {
		filters = "user_uid = ?"
		args = []interface{}{userUID[0]}
	}

	query := l.db.WithContext(ctx).Model(&auditLogs.AuditLog{})
	if filters != "" {
		query = query.Where(filters, args...)
	}

	err := query.Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	query = paginator(ctx, l.db, pagination)
	if filters != "" {
		query = query.Where(filters, args...)
	}

	if err = query.Model(&logs).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return &logs, totalRecords, nil
}
