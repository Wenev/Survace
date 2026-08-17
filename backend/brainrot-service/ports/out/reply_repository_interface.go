package out

import (
	"context"
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
)

type ReplyRepository interface {
	Create(ctx context.Context, reply *domain.Reply) error
	FindById(ctx context.Context, id int32) (*domain.Reply, error)
	FindByUserId(ctx context.Context, userId int32) ([]*domain.Reply, error)
	FindByCommentId(ctx context.Context, commentId int32) ([]*domain.Reply, error)
	CountByCommentId(ctx context.Context, commentId int32) (int64, error)
	DeleteUserComment(ctx context.Context, userId int32) error
	Delete(ctx context.Context, id int32) error
}
