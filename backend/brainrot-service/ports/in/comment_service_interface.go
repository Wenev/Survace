package in

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
)

type CommentService interface {
	AddComment(ctx context.Context, userId int32, videoId int32, text string) (int32, string, error)
	GetCommentsByVideo(ctx context.Context, videoId int32) (int32, string, []*domain.Comment, []int64, []int64, error)
	DeleteComment(ctx context.Context, commentId int32) (int32, string, error)
}
