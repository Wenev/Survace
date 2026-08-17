package in

import (
	"context"
)

type LikeReplyService interface {
	LikeReply(ctx context.Context, replyId int32, userId int32) (int32, string, error)
	UnlikeReply(ctx context.Context, replyId int32, userId int32) (int32, string, error)
	IsReplyLiked(ctx context.Context, replyId int32, userId int32) (bool, error)
}
