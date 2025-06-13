package repository

import (
	"context"
	"gorm.io/gorm"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/crypto"
	"ocserv-bakend/pkg/database"
	"time"
)

type UserRepository struct {
	db *gorm.DB
}
type UserRepositoryInterface interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	CreateToken(ctx context.Context, id uint, uid string, rememberMe bool, isAdmin bool) (string, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
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
