package repository

import (
	"context"
	"gorm.io/gorm"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/crypto"
	"ocserv-bakend/pkg/database"
	"ocserv-bakend/pkg/request"
	"time"
)

type UserRepository struct {
	db *gorm.DB
}
type UserRepositoryInterface interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByUID(ctx context.Context, uid string) (*models.User, error)
	CreateToken(ctx context.Context, id uint, uid string, rememberMe bool, isAdmin bool) (string, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	Users(ctx context.Context, pagination *request.Pagination) (*[]models.User, int64, error)
	ChangePassword(ctx context.Context, uid, password, salt string) error
	DeleteUser(ctx context.Context, uid string) error
	UpdateLastLogin(ctx context.Context, user *models.User) error
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.Get(),
	}
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateToken(ctx context.Context, id uint, uid string, rememberMe bool, isAdmin bool) (string, error) {
	expire := time.Now().Add(24 * time.Hour)
	if rememberMe {
		expire = expire.AddDate(0, 1, 0)
	}

	access, err := crypto.GenerateAccessToken(uid, expire.Unix(), isAdmin)
	if err != nil {
		return "", err
	}

	err = r.db.WithContext(ctx).Create(
		&models.UserToken{
			UserID:   id,
			Token:    access,
			ExpireAt: expire,
		},
	).Error
	if err != nil {
		return "", err
	}
	return access, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Users(ctx context.Context, pagination *request.Pagination) (*[]models.User, int64, error) {
	var totalRecords int64

	whereFilters := "is_admin = false"

	if err := r.db.WithContext(ctx).Model(&models.User{}).Where(whereFilters).Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	var staffs []models.User
	txPaginator := paginator(ctx, r.db, pagination)
	err := txPaginator.Model(&staffs).Where(whereFilters).Find(&staffs).Error
	if err != nil {
		return nil, 0, err
	}
	return &staffs, totalRecords, nil
}

func (r *UserRepository) ChangePassword(ctx context.Context, uid, password, salt string) error {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("uid = ?", uid).Updates(
		map[string]interface{}{
			"password": password,
			"salt":     salt,
		},
	).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, uid string) error {
	var user models.User
	err := r.db.WithContext(ctx).Where("uid = ? AND is_admin = ?", uid, false).First(&user).Error
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Delete(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByUID(ctx context.Context, uid string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("uid = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"last_login": user.LastLogin,
		}).Error
}
