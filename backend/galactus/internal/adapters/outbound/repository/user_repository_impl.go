package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/config"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/adapters/outbound/cache"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/app/domain"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
	"mime"
	"time"
)

type UserRepositoryImpl struct {
	db    *gorm.DB
	cache *cache.MemcachedConnection
	minio *config.MinIOClient
}

func NewUserRepository(db *gorm.DB, cache *cache.MemcachedConnection, minio *config.MinIOClient) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db:    db,
		cache: cache,
		minio: minio,
	}
}

func (u *UserRepositoryImpl) Create(ctx context.Context, user *domain.User) error {
	err := u.db.WithContext(ctx).Create(user).Error
	return err
}

func (u *UserRepositoryImpl) FindAll(ctx context.Context) ([]*domain.User, error) {
	var users []*domain.User
	result := u.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return users, nil
}

func (u *UserRepositoryImpl) FindByID(ctx context.Context, id int32) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).First(&user, id)
	log.Print("bring")
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	log.Printf("hello %d", user.ID)
	return &user, nil
}

func (u *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImpl) FindByEmailOrUsername(ctx context.Context, emailOrUsername string) (*domain.User, error) {
	var user domain.User
	result := u.db.WithContext(ctx).Where("username = ?", emailOrUsername).Or("email = ?", emailOrUsername).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (u *UserRepositoryImpl) Update(ctx context.Context, user *domain.User) error {
	err := u.db.WithContext(ctx).Save(&user).Error
	return err
}

func (u *UserRepositoryImpl) Delete(ctx context.Context, userID int32) error {
	return u.db.WithContext(ctx).Delete(&domain.User{}, userID).Error
}

func (u *UserRepositoryImpl) StoreCache(key string, data interface{}, expirationTime time.Duration) error {
	err := u.cache.Set(key, data, expirationTime)
	log.Print(err)
	if err != nil {
		return status.Error(codes.Internal, "Internal Server Error")
	}
	return nil
}

func (u *UserRepositoryImpl) GetCache(key string, dest interface{}) error {
	err := u.cache.Get(key, dest)
	log.Print(err)
	if err != nil {
		return status.Error(codes.NotFound, "Cache key not found")
	}
	return nil
}

func (u *UserRepositoryImpl) DeleteCache(key string) error {
	err := u.cache.Delete(key)
	if err != nil {
		return status.Error(codes.Internal, "Internal Server Error")
	}
	return nil
}

func (u *UserRepositoryImpl) UploadAvatar(ctx context.Context, userID int32, avatarBlob []byte) (string, error) {
	if u.minio == nil {
		return "", status.Error(codes.Internal, "MinIO client not initialized")
	}
	objectName := fmt.Sprintf("avatar_%d_%d.jpg", userID, time.Now().UnixNano())
	contentType := mime.TypeByExtension(".jpg")
	reader := bytes.NewReader(avatarBlob)
	_, err := u.minio.Client.PutObject(ctx, u.minio.Bucket, objectName, reader, int64(len(avatarBlob)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", status.Error(codes.Internal, "Failed to upload avatar to MinIO")
	}
	url := fmt.Sprintf("http://%s/%s/%s", u.minio.Client.EndpointURL().Host, u.minio.Bucket, objectName)
	return url, nil
}
