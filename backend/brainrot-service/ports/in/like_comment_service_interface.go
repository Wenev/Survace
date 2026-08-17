package in

import (
	"context"
)

type LikeCommentService interface {
	LikeComment(ctx context.Context, userId int32, commentId int32) (int32, string, error)
	UnlikeComment(ctx context.Context, userId int32, commentId int32) (int32, string, error)
	IsCommentLiked(ctx context.Context, userId int32, commentId int32) (bool, error)
}
