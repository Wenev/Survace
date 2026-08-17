package in

import (
	"context"
	"github.com/Wenev/Survace/galactus/internal/app/domain"
	"time"
)

type UserService interface {
	Register(ctx context.Context, username, email, password, otpCode string, dateOfBirth time.Time) (int32, string,
		*domain.User,
		error)
	GoogleRegister(ctx context.Context, email, username string) (int32, string,
		*domain.User,
		error)
	Login(ctx context.Context, emailOrUsername, password string) (int32, string, string, *domain.User, error)
	GoogleLogin(ctx context.Context, email string) (int32, string, string, *domain.User, error)
	UpdateUser(ctx context.Context, userID int32, username, email, password string, avatarBlob []byte) (int32, string, *domain.User, error)
	SendOTP(ctx context.Context, email string) (int32, string, error)
	FindById(ctx context.Context, userID int32) (int32, string, *domain.User, error)
	FindByUsername(ctx context.Context, username string) (int32, string, *domain.User, error)
	IsValid(ctx context.Context, userID int32) (int32, string, error)
	ForgetPassword(ctx context.Context, email string) (int32, string, error)
	ResetPassword(ctx context.Context, email, otpCode, newPassword string) (int32, string, error)
}
