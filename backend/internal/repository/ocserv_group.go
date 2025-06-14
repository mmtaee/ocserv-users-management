package repository

import (
	"context"
	"gorm.io/gorm"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/database"
	"ocserv-bakend/pkg/request"
)

type OcservGroupRepository struct {
	db *gorm.DB
}

type OcservGroupRepositoryInterface interface {
	Groups(ctx context.Context, pagination *request.Pagination) (*[]models.OcservGroup, int64, error)
}

func NewOcservGroupRepository() *OcservGroupRepository {
	return &OcservGroupRepository{db: database.Get()}
}

func (o *OcservGroupRepository) Groups(ctx context.Context, pagination *request.Pagination) (*[]models.OcservGroup, int64, error) {
	var totalRecords int64

	err := o.db.WithContext(ctx).Model(&models.OcservGroup{}).Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	var ocservGroups []models.OcservGroup
	txPaginator := paginator(ctx, o.db, pagination)
	err = txPaginator.Model(&ocservGroups).Find(&ocservGroups).Error
	if err != nil {
		return nil, 0, err
	}
	return &ocservGroups, totalRecords, nil

}
