package repository

import (
	"context"
	"gorm.io/gorm"
	"log"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/config"
	"ocserv-bakend/pkg/database"
	ocApi "ocserv-bakend/pkg/oc_api"
	"ocserv-bakend/pkg/request"
)

type OcservGroupRepository struct {
	db    *gorm.DB
	ocApi ocApi.OcGroupApiRepositoryInterface
}

type OcservGroupRepositoryInterface interface {
	Groups(ctx context.Context, pagination *request.Pagination) (*[]models.OcservGroup, int64, error)
	GroupsLookup(ctx context.Context) ([]string, error)
	GetByID(ctx context.Context, id string) (*models.OcservGroup, error)
	Create(ctx context.Context, ocservGroup *models.OcservGroup) (*models.OcservGroup, error)
	Update(ctx context.Context, ocservGroup *models.OcservGroup) (*models.OcservGroup, error)
	Delete(ctx context.Context, id string) (*models.OcservGroup, error)
}

func NewOcservGroupRepository() *OcservGroupRepository {
	apiURLService := config.Get().APIURLService
	return &OcservGroupRepository{
		db:    database.Get(),
		ocApi: ocApi.NewOcGroupApiRepository(apiURLService),
	}
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

func (o *OcservGroupRepository) GroupsLookup(ctx context.Context) ([]string, error) {
	var ocservGroups []models.OcservGroup

	err := o.db.WithContext(ctx).Model(&models.OcservGroup{}).Select("name").Find(&ocservGroups).Error
	if err != nil {
		return nil, err
	}

	groups := make([]string, 0, len(ocservGroups))
	for _, ocservGroup := range ocservGroups {
		groups = append(groups, ocservGroup.Name)
	}
	return groups, nil
}

func (o *OcservGroupRepository) GetByID(ctx context.Context, id string) (*models.OcservGroup, error) {
	var ocservGroup models.OcservGroup
	err := o.db.WithContext(ctx).Where("id = ?", id).Find(&ocservGroup).Error
	if err != nil {
		return nil, err
	}
	return &ocservGroup, nil
}

func (o *OcservGroupRepository) Create(ctx context.Context, ocservGroup *models.OcservGroup) (*models.OcservGroup, error) {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(ocservGroup).Error; err != nil {
			return err
		}
		if err := o.ocApi.CreateGroupApi(ctx, ocservGroup.Name, ocservGroup.Config); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return ocservGroup, nil
}

func (o *OcservGroupRepository) Update(ctx context.Context, ocservGroup *models.OcservGroup) (*models.OcservGroup, error) {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(ocservGroup).Save(ocservGroup).Error; err != nil {
			return err
		}
		if err := o.ocApi.CreateGroupApi(ctx, ocservGroup.Name, ocservGroup.Config); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ocservGroup, nil
}

func (o *OcservGroupRepository) Delete(ctx context.Context, id string) (*models.OcservGroup, error) {
	var group models.OcservGroup
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&group).Error; err != nil {
			return err
		}

		if err := tx.Delete(&group).Error; err != nil {
			log.Println("err: ", err)
			return err
		}

		if err := o.ocApi.DeleteGroupApi(ctx, group.Name); err != nil {
			return err
		}
		return nil
	})
	return &group, err
}
