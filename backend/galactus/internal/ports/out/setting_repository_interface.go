package out

import (
	"context"
	"github.com/Wenev/Survace/galactus/internal/app/domain"
)

type SettingRepository interface {
	Create(ctx context.Context, setting *domain.Setting) error
	FindByUserID(ctx context.Context, userID int32) (*domain.Setting, error)
	Edit(ctx context.Context, setting *domain.Setting) error
	Delete(ctx context.Context, userID int32) error
}
