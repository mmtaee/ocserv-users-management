package repository

import (
	"context"
	"gorm.io/gorm"
	auditLogs "ocserv-bakend/pkg/audit_log"
	"ocserv-bakend/pkg/database"
	"ocserv-bakend/pkg/request"
)

type LogsRepository struct {
	db *gorm.DB
}
type LogsRepositoryInterface interface {
	UsersLogs(ctx context.Context, pagination *request.Pagination, userUID string) (*[]auditLogs.AuditLog, int64, error)
	Logs(ctx context.Context, pagination *request.Pagination) (*[]auditLogs.AuditLog, int64, error)
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

func (l *LogsRepository) Logs(ctx context.Context, pagination *request.Pagination) (*[]auditLogs.AuditLog, int64, error) {
	var (
		totalRecords int64
		logs         []auditLogs.AuditLog
	)
	if err := l.db.WithContext(ctx).Model(&auditLogs.AuditLog{}).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	txPaginator := paginator(ctx, l.db, pagination)
	if err := txPaginator.Model(&logs).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return &logs, totalRecords, nil
}
