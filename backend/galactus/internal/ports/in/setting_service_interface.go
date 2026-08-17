package in

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/app/domain"
)

type SettingService interface {
	CreateSetting(ctx context.Context, userID int32) (int32, string, error)
	EnableNewFollowerNotification(ctx context.Context, userID int32, enable bool) (int32, string, error)
	EnableMessageNotification(ctx context.Context, userID int32, enable bool) (int32, string, error)
	EnableMentionsNotification(ctx context.Context, userID int32, enable bool) (int32, string, error)
	EnableLikeTabVisibility(ctx context.Context, userID int32, enable bool) (int32, string, error)
	EnablePrivateAccount(ctx context.Context, userID int32, enable bool) (int32, string, error)
	EditChatRestriction(ctx context.Context, userID int32, restriction domain.ChatRestrictionType) (int32, string, error)
	DeleteAccount(ctx context.Context, userID int32) (int32, string, error)
	FindUserSetting(ctx context.Context, userID int32) (*domain.Setting, int32, string, error)
}
