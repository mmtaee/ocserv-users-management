package repository

import (
	"context"
	"gorm.io/gorm"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/database"
)

type SystemRepository struct {
	db *gorm.DB
}

type SystemRepositoryInterface interface {
	System(ctx context.Context) (*models.System, error)
	SystemUpdate(ctx context.Context, system *models.System) (*models.System, error)
}

func NewSystemRepository() *SystemRepository {
	return &SystemRepository{
		db: database.Get(),
	}
}

func (s *SystemRepository) System(ctx context.Context) (*models.System, error) {
	var system models.System
	err := s.db.WithContext(ctx).First(&system).Error
	if err != nil {
		return nil, err
	}
	return &system, nil
}

func (s *SystemRepository) SystemUpdate(ctx context.Context, system *models.System) (*models.System, error) {
	err := s.db.WithContext(ctx).Save(&system).Error
	if err != nil {
		return nil, err
	}
	return system, nil
}
