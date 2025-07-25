package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"ocserv-bakend/pkg/request"
)

func paginator(ctx context.Context, db *gorm.DB, pagination *request.Pagination) *gorm.DB {
	if pagination.Order == "" {
		pagination.Order = "id"
		pagination.Sort = "ASC"
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	
	order := fmt.Sprintf("%s %s", pagination.Order, pagination.Sort)

	return db.WithContext(ctx).Order(order).Limit(pagination.PageSize).Offset(offset)
}
