package grpc

import (
	"context"
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/app/domain"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/ports/in"
	"github.com/Acad600-TPA/WEB-WE-251/proto/gen/controller"
	"github.com/Acad600-TPA/WEB-WE-251/proto/gen/dto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

type GrpcServer struct {
	grpcPort       int
	server         *grpc.Server
	userService    in.UserService
	settingService in.SettingService
	controller.GalactusControllerServer
}

func NewGrpcServer(grpcPort int, userService in.UserService, settingService in.SettingService) *GrpcServer {
	return &GrpcServer{
		grpcPort:       grpcPort,
		server:         grpc.NewServer(),
		userService:    userService,
		settingService: settingService,
	}
}

func (g *GrpcServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.grpcPort))
	if err != nil {
		fmt.Print("HELP")
		return err
	}
	controller.RegisterGalactusControllerServer(g.server, g)
	err = g.server.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (g *GrpcServer) Stop() {
	g.server.GracefulStop()
}

func (g *GrpcServer) Register(ctx context.Context, in *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	code, message, user, err := g.userService.Register(ctx, in.Username, in.Email, in.Password, in.OTPCode,
		in.DateOfBirth.AsTime())
	if err != nil {
		return nil, err
	}

	var userId int32 = 0
	if user != nil {
		userId = user.ID
		_, _, setErr := g.settingService.CreateSetting(ctx, userId)
		if setErr != nil {
			return nil, setErr
		}
	}

	resp := &dto.RegisterResponse{
		StatusCode: code,
		Message:    message,
		UserId:     userId,
	}
	return resp, nil
}

func (g *GrpcServer) GoogleRegister(ctx context.Context, in *dto.GoogleRegisterRequest) (*dto.RegisterResponse, error) {
	code, message, user, err := g.userService.GoogleRegister(ctx, in.Email, in.Username)
	if err != nil {
		return nil, err
	}

	var userId int32 = 0
	if user != nil {
		userId = user.ID
		_, _, setErr := g.settingService.CreateSetting(ctx, userId)
		if setErr != nil {
			return nil, setErr
		}
	}

	resp := &dto.RegisterResponse{
		StatusCode: code,
		Message:    message,
		UserId:     userId,
	}
	return resp, nil
}

func (g *GrpcServer) Login(ctx context.Context, in *dto.LoginRequest) (*dto.LoginResponse, error) {
	code, message, token, user, err := g.userService.Login(ctx, in.EmailOrUsername, in.Password)
	if err != nil {
		return nil, err
	}

	var userId int32
	if user != nil {
		userId = user.ID
	}

	resp := &dto.LoginResponse{
		StatusCode: code,
		Message:    message,
		Token:      token,
		UserId:     userId,
	}
	return resp, nil
}
func (g *GrpcServer) GoogleLogin(ctx context.Context, in *dto.GoogleLoginRequest) (*dto.LoginResponse, error) {
	code, message, token, user, err := g.userService.GoogleLogin(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	var userId int32
	if user != nil {
		userId = user.ID
	}

	resp := &dto.LoginResponse{
		StatusCode: code,
		Message:    message,
		Token:      token,
		UserId:     userId,
	}
	return resp, nil
}
func (g *GrpcServer) UpdateUser(ctx context.Context, in *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	code, message, user, err := g.userService.UpdateUser(ctx, in.UserId, in.Username, in.Email, in.Password, in.AvatarBlob)
	if err != nil {
		return &dto.UpdateUserResponse{
			StatusCode: code,
			UserId:     in.UserId,
			Message:    message,
			AvatarUrl:  "",
		}, err
	}
	resp := &dto.UpdateUserResponse{
		StatusCode: code,
		UserId:     user.ID,
		Message:    message,
		AvatarUrl:  user.AvatarURL,
	}
	return resp, nil
}
func (g *GrpcServer) SendOTP(ctx context.Context, in *dto.SendOTPRequest) (*dto.SendOTPResponse, error) {
	log.Printf("SendOTP called with emailwaa: %s", in.Email)
	code, message, err := g.userService.SendOTP(ctx, in.Email)
	if err != nil {
		return nil, err
	}
	resp := &dto.SendOTPResponse{
		StatusCode: code,
		Message:    message,
	}
	return resp, nil
}
func (g *GrpcServer) FindByUserId(ctx context.Context, in *dto.FindByUserIdRequest) (*dto.FindByUserIdResponse, error) {
	code, message, user, err := g.userService.FindById(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	userData := &dto.User{
		UserId:        user.ID,
		Username:      user.Username,
		Email:         user.Email,
		DateOfBirth:   timestamppb.New(user.DateOfBirth),
		EmailVerified: user.IsEmailVerified,
	}

	resp := &dto.FindByUserIdResponse{
		StatusCode: code,
		Message:    message,
		User:       userData,
	}
	return resp, nil
}

func (g *GrpcServer) FindByUsername(ctx context.Context, in *dto.FindByUsernameRequest) (*dto.FindByUsernameResponse, error) {
	code, message, user, err := g.userService.FindByUsername(ctx, in.Username)
	if err != nil {
		return nil, err
	}

	var userData *dto.User
	if user != nil {
		userData = &dto.User{
			UserId:        user.ID,
			Username:      user.Username,
			Email:         user.Email,
			DateOfBirth:   timestamppb.New(user.DateOfBirth),
			EmailVerified: user.IsEmailVerified,
			AvatarUrl:     user.AvatarURL,
		}
	}

	resp := &dto.FindByUsernameResponse{
		StatusCode: code,
		Message:    message,
		User:       userData,
	}
	return resp, nil
}

func (g *GrpcServer) EnableNewFollowerNotification(ctx context.Context, in *dto.EnableSettingRequest) (*dto.SettingResponse, error) {
	code, message, err := g.settingService.EnableNewFollowerNotification(ctx, in.UserId, in.Enable)
	if err != nil {
		return nil, err
	}
	return &dto.SettingResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) EnableMessageNotification(ctx context.Context, in *dto.EnableSettingRequest) (*dto.SettingResponse, error) {
	code, message, err := g.settingService.EnableMessageNotification(ctx, in.UserId, in.Enable)
	if err != nil {
		return nil, err
	}
	return &dto.SettingResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) EnableMentionsNotification(ctx context.Context, in *dto.EnableSettingRequest) (*dto.SettingResponse, error) {
	code, message, err := g.settingService.EnableMentionsNotification(ctx, in.UserId, in.Enable)
	if err != nil {
		return nil, err
	}
	return &dto.SettingResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) EnableLikeTabVisibility(ctx context.Context, in *dto.EnableSettingRequest) (*dto.SettingResponse, error) {
	code, message, err := g.settingService.EnableLikeTabVisibility(ctx, in.UserId, in.Enable)
	if err != nil {
		return nil, err
	}
	return &dto.SettingResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) EnablePrivateAccount(ctx context.Context, in *dto.EnableSettingRequest) (*dto.SettingResponse, error) {
	code, message, err := g.settingService.EnablePrivateAccount(ctx, in.UserId, in.Enable)
	if err != nil {
		return nil, err
	}
	return &dto.SettingResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) EditChatRestriction(ctx context.Context, in *dto.EditChatRestrictionRequest) (*dto.SettingResponse, error) {
	code, message, err := g.settingService.EditChatRestriction(ctx, in.UserId, domain.ChatRestrictionType(in.ChatRestriction))
	if err != nil {
		return nil, err
	}
	return &dto.SettingResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) DeleteAccount(ctx context.Context, in *dto.DeleteAccountRequest) (*dto.SettingResponse, error) {
	code, message, err := g.settingService.DeleteAccount(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &dto.SettingResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) IsValid(ctx context.Context, in *dto.IsValidRequest) (*dto.IsValidResponse, error) {
	code, message, err := g.userService.IsValid(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &dto.IsValidResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) FindUserSetting(ctx context.Context, in *dto.FindUserSettingRequest) (*dto.FindUserSettingResponse, error) {
	setting, code, message, err := g.settingService.FindUserSetting(ctx, in.Id)
	if err != nil {
		return &dto.FindUserSettingResponse{
			StatusCode: code,
			Message:    message,
			Setting:    nil,
		}, err
	}
	if setting == nil {
		return &dto.FindUserSettingResponse{
			StatusCode: code,
			Message:    message,
			Setting:    nil,
		}, nil
	}
	settingProto := &dto.Setting{
		Id:                      setting.ID,
		UserId:                  setting.UserID,
		ChatRestriction:         string(setting.ChatRestriction),
		Private:                 setting.Private,
		NewFollowerNotification: setting.NewFollowerNotification,
		MessageNotification:     setting.MessageNotification,
		MentionsNotification:    setting.MentionsNotification,
		LikeTabVisibility:       setting.LikeTabVisibility,
	}
	return &dto.FindUserSettingResponse{
		StatusCode: code,
		Message:    message,
		Setting:    settingProto,
	}, nil
}

func (g *GrpcServer) ForgetPassword(ctx context.Context, in *dto.ForgetPasswordRequest) (*dto.ForgetPasswordResponse, error) {
	code, message, err := g.userService.ForgetPassword(ctx, in.Email)
	if err != nil {
		return &dto.ForgetPasswordResponse{
			StatusCode: code,
			Message:    message,
		}, err
	}
	return &dto.ForgetPasswordResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}

func (g *GrpcServer) ResetPassword(ctx context.Context, in *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, error) {
	code, message, err := g.userService.ResetPassword(ctx, in.Email, in.OtpCode, in.NewPassword)
	if err != nil {
		return &dto.ResetPasswordResponse{
			StatusCode: code,
			Message:    message,
		}, err
	}
	return &dto.ResetPasswordResponse{
		StatusCode: code,
		Message:    message,
	}, nil
}
