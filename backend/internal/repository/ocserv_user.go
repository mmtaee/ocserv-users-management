package repository

import (
	"context"
	"gorm.io/gorm"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/config"
	"ocserv-bakend/pkg/database"
	ocApi "ocserv-bakend/pkg/oc_api"
	"ocserv-bakend/pkg/request"
	"time"
)

type OcservUserRepository struct {
	db    *gorm.DB
	ocApi ocApi.OcUserApiRepositoryInterface
}

type OcservUserRepositoryInterface interface {
	Users(ctx context.Context, pagination *request.Pagination) (*[]models.OcservUser, int64, error)
	Create(ctx context.Context, user *models.OcservUser) (*models.OcservUser, error)
	GetByUID(ctx context.Context, uid string) (*models.OcservUser, error)
	Update(ctx context.Context, ocservUser *models.OcservUser) (*models.OcservUser, error)
	Lock(ctx context.Context, uid string) error
	UnLock(ctx context.Context, uid string) error
	Delete(ctx context.Context, uid string) error
	TenDaysStats(ctx context.Context) (*[]models.DailyTraffic, error)
}

func NewtOcservUserRepository() *OcservUserRepository {
	apiURLService := config.Get().APIURLService
	return &OcservUserRepository{
		db:    database.Get(),
		ocApi: ocApi.NewOcUserApiRepository(apiURLService),
	}
}

func (o *OcservUserRepository) Users(ctx context.Context, pagination *request.Pagination) (*[]models.OcservUser, int64, error) {
	var totalRecords int64

	err := o.db.WithContext(ctx).Model(&models.OcservUser{}).Count(&totalRecords).Error
	if err != nil {
		return nil, 0, err
	}

	var ocservUser []models.OcservUser
	txPaginator := paginator(ctx, o.db, pagination)
	err = txPaginator.Model(&ocservUser).Find(&ocservUser).Error
	if err != nil {
		return nil, 0, err
	}
	return &ocservUser, totalRecords, nil
}

func (o *OcservUserRepository) Create(ctx context.Context, user *models.OcservUser) (*models.OcservUser, error) {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if err := o.ocApi.CreateUserApi(ctx, user.Group, user.Username, user.Password); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return user, err
}

func (o *OcservUserRepository) GetByUID(ctx context.Context, uid string) (*models.OcservUser, error) {
	var ocservUser models.OcservUser
	err := o.db.WithContext(ctx).Where("uid = ?", uid).First(&ocservUser).Error
	if err != nil {
		return nil, err
	}
	return &ocservUser, nil
}

func (o *OcservUserRepository) Update(ctx context.Context, ocservUser *models.OcservUser) (*models.OcservUser, error) {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&ocservUser).Error; err != nil {
			return err
		}
		if err := o.ocApi.CreateUserApi(ctx, ocservUser.Group, ocservUser.Username, ocservUser.Password); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ocservUser, nil
}

func (o *OcservUserRepository) Lock(ctx context.Context, uid string) error {
	var ocservUser models.OcservUser
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("uid = ?", uid).First(&ocservUser).Error; err != nil {
			return err
		}
		if err := tx.Updates(map[string]interface{}{"is_locked": true}).Error; err != nil {
			return err
		}

		if err := o.ocApi.LockUserApi(ctx, ocservUser.Username); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (o *OcservUserRepository) UnLock(ctx context.Context, uid string) error {
	var ocservUser models.OcservUser
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("uid = ?", uid).First(&ocservUser).Error; err != nil {
			return err
		}
		if err := tx.Updates(map[string]interface{}{"is_locked": false}).Error; err != nil {
			return err
		}

		if err := o.ocApi.UnLockUserApi(ctx, ocservUser.Username); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (o *OcservUserRepository) Delete(ctx context.Context, uid string) error {
	var ocservUser models.OcservUser
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("uid = ?", uid).First(&ocservUser).Error; err != nil {
			return err
		}
		if err := tx.Delete(&ocservUser).Error; err != nil {
			return err
		}
		if err := o.ocApi.DeleteUserApi(ctx, ocservUser.Username); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (o *OcservUserRepository) TenDaysStats(ctx context.Context) (*[]models.DailyTraffic, error) {
	var results []models.DailyTraffic

	err := o.db.WithContext(ctx).
		Model(&models.OcservUserTrafficStatistics{}).
		Select(`
		DATE(created_at) AS date,
		SUM(rx) / 1073741824.0 AS rx,
		SUM(tx) / 1073741824.0 AS tx`).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -10)).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	return &results, nil
}
