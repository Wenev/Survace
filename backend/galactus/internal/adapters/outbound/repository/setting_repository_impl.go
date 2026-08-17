package repository

import (
	"context"
	"github.com/Wenev/Survace/galactus/internal/app/domain"
	"gorm.io/gorm"
)

type SettingRepositoryImpl struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepositoryImpl {
	return &SettingRepositoryImpl{db: db}
}

func (r *SettingRepositoryImpl) Create(ctx context.Context, setting *domain.Setting) error {
	return r.db.WithContext(ctx).Create(setting).Error
}

func (r *SettingRepositoryImpl) FindByUserID(ctx context.Context, userID int32) (*domain.Setting, error) {
	var setting domain.Setting
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

func (r *SettingRepositoryImpl) Edit(ctx context.Context, setting *domain.Setting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *SettingRepositoryImpl) Delete(ctx context.Context, userID int32) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.Setting{}).Error
}
