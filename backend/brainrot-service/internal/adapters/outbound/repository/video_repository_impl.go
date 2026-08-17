package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Wenev/Survace/brainrot-service/config"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log"
)

type VideoRepositoryImpl struct {
	minio *config.MinIOClient
	db    *gorm.DB
}

func NewVideoRepository(db *gorm.DB, minio *config.MinIOClient) *VideoRepositoryImpl {
	return &VideoRepositoryImpl{
		minio: minio,
		db:    db,
	}
}

func (v *VideoRepositoryImpl) Create(ctx context.Context, video *domain.Video) error {
	err := v.db.WithContext(ctx).Create(video).Error
	return err
}

func (v *VideoRepositoryImpl) FindById(ctx context.Context, id int32) (*domain.Video, error) {
	var video domain.Video
	result := v.db.WithContext(ctx).First(&video, id)
	log.Print("bring")
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &video, nil
}

func (v *VideoRepositoryImpl) FindByUserId(ctx context.Context, userId int32) ([]*domain.Video, error) {
	var video []*domain.Video
	err := v.db.WithContext(ctx).Where("user_id = ? AND is_draft = ?", userId, false).Find(&video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (v *VideoRepositoryImpl) FindByUserIdAndAll(ctx context.Context, userId int32) ([]*domain.Video, error) {
	var video []*domain.Video
	err := v.db.WithContext(ctx).Where("user_id = ?", userId).Find(&video).Error
	if err != nil {
		return nil, err
	}
	return video, nil
}

func (v *VideoRepositoryImpl) Update(ctx context.Context, video *domain.Video) error {
	err := v.db.WithContext(ctx).Save(&video).Error
	return err
}

func (v *VideoRepositoryImpl) UploadMinIO(ctx context.Context, objectName string, data []byte,
	contentType string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := v.minio.Client.PutObject(ctx, v.minio.Bucket, objectName, reader, int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", status.Errorf(codes.Internal, "Failed to upload video")
	}

	url := fmt.Sprintf("http://localhost:9000/%s/%s", v.minio.Bucket, objectName)
	return url, nil
}

func (v *VideoRepositoryImpl) Delete(ctx context.Context, id int32) error {
	result := v.db.WithContext(ctx).Delete(&domain.Video{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (v *VideoRepositoryImpl) DeleteByUserId(ctx context.Context, userId int32) error {
	err := v.db.WithContext(ctx).Where("user_id = ?", userId).Delete(&domain.Video{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (v *VideoRepositoryImpl) FindNewAndPopular(ctx context.Context, limit int, offset int) ([]*domain.Video, error) {
	var videos []*domain.Video
	err := v.db.WithContext(ctx).
		Where("visibility = ? AND is_draft = ?", "public", false).
		Order("created_at DESC").
		Order("view_count DESC").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (v *VideoRepositoryImpl) FindWatchedVideoIDs(ctx context.Context, userId int32) ([]int32, error) {
	var watched []struct{ VideoID int32 }
	err := v.db.WithContext(ctx).
		Table("video_views").
		Select("video_id").
		Where("user_id = ?", userId).
		Scan(&watched).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int32, 0, len(watched))
	for _, w := range watched {
		ids = append(ids, w.VideoID)
	}
	return ids, nil
}

func (v *VideoRepositoryImpl) FindRandomVideos(ctx context.Context, limit int, offset int) ([]*domain.Video, error) {
	var videos []*domain.Video
	err := v.db.WithContext(ctx).
		Where("is_draft = ?", false).
		Order("RANDOM()").
		Limit(limit).
		Offset(offset).
		Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
