package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/app/domain"
	"time"
)

type UserRespository interface {
	Create(ctx context.Context, user *domain.User) error
	FindAll(ctx context.Context) ([]*domain.User, error)
	FindByID(ctx context.Context, id int32) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByEmailOrUsername(ctx context.Context, emailOrUsername string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, userID int32) error
	StoreCache(key string, code interface{}, expirationTime time.Duration) error
	GetCache(key string, dest interface{}) error
	UploadAvatar(ctx context.Context, userID int32, avatarBlob []byte) (string, error)
	DeleteCache(key string) error
}
