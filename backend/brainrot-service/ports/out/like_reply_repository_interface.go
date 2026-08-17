package out

import (
	"context"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
)

type LikeReplyRepository interface {
	Create(ctx context.Context, like *domain.LikeReply) error
	FindByReplyIdAndUserId(ctx context.Context, replyId int32, userId int32) (*domain.LikeReply, error)
	Delete(ctx context.Context, id int32) error
	DeleteByReplyIdAndUserId(ctx context.Context, replyId int32, userId int32) error
	CountByReplyId(ctx context.Context, replyId int32) (int64, error)
}
