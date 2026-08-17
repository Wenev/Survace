package app

import (
	"context"
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/adapters/outbound/cache"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/app/domain"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/ports/out"
	"google.golang.org/grpc/codes"
	"time"
)

type SettingServiceImpl struct {
	repo     out.SettingRepository
	userRepo out.UserRespository
	cache    *cache.MemcachedConnection
}

func NewSettingService(repo out.SettingRepository, userRepo out.UserRespository, cache *cache.MemcachedConnection) *SettingServiceImpl {
	return &SettingServiceImpl{repo: repo, userRepo: userRepo, cache: cache}
}

func (s *SettingServiceImpl) CreateSetting(ctx context.Context, userID int32) (int32, string, error) {
	setting := &domain.Setting{UserID: userID}
	err := s.repo.Create(ctx, setting)
	if err != nil {
		return int32(codes.Internal), "Failed to create setting", err
	}
	_ = s.cache.Set(s.settingCacheKey(userID), setting, 10*time.Minute)
	return int32(codes.OK), "Setting created successfully", nil
}

func (s *SettingServiceImpl) EnableNewFollowerNotification(ctx context.Context, userID int32, enable bool) (int32, string, error) {
	setting, err := s.getSettingWithCache(ctx, userID)
	if err != nil || setting == nil {
		return int32(codes.NotFound), "Setting not found", err
	}
	setting.NewFollowerNotification = enable
	err = s.repo.Edit(ctx, setting)
	if err != nil {
		return int32(codes.Internal), "Failed to update notification", err
	}
	_ = s.cache.Set(s.settingCacheKey(userID), setting, 10*time.Minute)
	return int32(codes.OK), "New follower notification updated", nil
}

func (s *SettingServiceImpl) EnableMessageNotification(ctx context.Context, userID int32, enable bool) (int32, string, error) {
	setting, err := s.getSettingWithCache(ctx, userID)
	if err != nil || setting == nil {
		return int32(codes.NotFound), "Setting not found", err
	}
	setting.MessageNotification = enable
	err = s.repo.Edit(ctx, setting)
	if err != nil {
		return int32(codes.Internal), "Failed to update message notification", err
	}
	_ = s.cache.Set(s.settingCacheKey(userID), setting, 10*time.Minute)
	return int32(codes.OK), "Message notification updated", nil
}

func (s *SettingServiceImpl) EnableMentionsNotification(ctx context.Context, userID int32, enable bool) (int32, string, error) {
	setting, err := s.getSettingWithCache(ctx, userID)
	if err != nil || setting == nil {
		return int32(codes.NotFound), "Setting not found", err
	}
	setting.MentionsNotification = enable
	err = s.repo.Edit(ctx, setting)
	if err != nil {
		return int32(codes.Internal), "Failed to update mentions notification", err
	}
	_ = s.cache.Set(s.settingCacheKey(userID), setting, 10*time.Minute)
	return int32(codes.OK), "Mentions notification updated", nil
}

func (s *SettingServiceImpl) EnableLikeTabVisibility(ctx context.Context, userID int32, enable bool) (int32, string, error) {
	setting, err := s.getSettingWithCache(ctx, userID)
	if err != nil || setting == nil {
		return int32(codes.NotFound), "Setting not found", err
	}
	setting.LikeTabVisibility = enable
	err = s.repo.Edit(ctx, setting)
	if err != nil {
		return int32(codes.Internal), "Failed to update like tab visibility", err
	}
	_ = s.cache.Set(s.settingCacheKey(userID), setting, 10*time.Minute)
	return int32(codes.OK), "Like tab visibility updated", nil
}

func (s *SettingServiceImpl) EnablePrivateAccount(ctx context.Context, userID int32, enable bool) (int32, string, error) {
	setting, err := s.getSettingWithCache(ctx, userID)
	if err != nil || setting == nil {
		return int32(codes.NotFound), "Setting not found", err
	}
	setting.Private = enable
	err = s.repo.Edit(ctx, setting)
	if err != nil {
		return int32(codes.Internal), "Failed to update private account", err
	}
	_ = s.cache.Set(s.settingCacheKey(userID), setting, 10*time.Minute)
	return int32(codes.OK), "Private account setting updated", nil
}

func (s *SettingServiceImpl) EditChatRestriction(ctx context.Context, userID int32, restriction domain.ChatRestrictionType) (int32, string, error) {
	setting, err := s.getSettingWithCache(ctx, userID)
	if err != nil || setting == nil {
		return int32(codes.NotFound), "Setting not found", err
	}
	setting.ChatRestriction = restriction
	err = s.repo.Edit(ctx, setting)
	if err != nil {
		return int32(codes.Internal), "Failed to update chat restriction", err
	}
	_ = s.cache.Set(s.settingCacheKey(userID), setting, 10*time.Minute)
	return int32(codes.OK), "Chat restriction updated", nil
}

func (s *SettingServiceImpl) DeleteAccount(ctx context.Context, userID int32) (int32, string, error) {
	_ = s.cache.Delete(s.settingCacheKey(userID))
	if err := s.repo.Delete(ctx, userID); err != nil {
		return int32(codes.Internal), "Failed to delete setting", err
	}
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return int32(codes.Internal), "Failed to delete user", err
	}
	return int32(codes.OK), "Account deleted successfully", nil
}

func (s *SettingServiceImpl) FindUserSetting(ctx context.Context, userID int32) (*domain.Setting, int32, string, error) {
	setting, err := s.getSettingWithCache(ctx, userID)
	if err != nil {
		return nil, int32(codes.NotFound), "Setting not found", err
	}
	if setting == nil {
		return nil, int32(codes.NotFound), "Setting not found", nil
	}
	return setting, int32(codes.OK), "Setting found", nil
}

func (s *SettingServiceImpl) getSettingWithCache(ctx context.Context, userID int32) (*domain.Setting, error) {
	var setting domain.Setting
	err := s.cache.Get(s.settingCacheKey(userID), &setting)
	if err == nil {
		return &setting, nil
	}
	settingPtr, err := s.repo.FindByUserID(ctx, userID)
	if err != nil || settingPtr == nil {
		return nil, err
	}
	_ = s.cache.Set(s.settingCacheKey(userID), settingPtr, 10*time.Minute)
	return settingPtr, nil
}

func (s *SettingServiceImpl) settingCacheKey(userID int32) string {
	return "setting:user:" + fmt.Sprint(userID)
}
