package in

import (
	"context"
	"github.com/Wenev/Survace/brainrot-service/internal/app/domain"
)

type ReplyService interface {
	AddReply(ctx context.Context, userId int32, commentId int32, text string) (int32, string, error)
	GetRepliesByComment(ctx context.Context, commentId int32) (int32, string, []*domain.Reply, []int64, error)
	DeleteReply(ctx context.Context, replyId int32) (int32, string, error)
}
