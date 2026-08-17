package app

import (
	"context"
	"fmt"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/app/domain"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/app/helper"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/ports/out"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type UserServiceImpl struct {
	userRepo out.UserRespository
}

func NewUserService(userRepo out.UserRespository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (u *UserServiceImpl) Register(ctx context.Context, username, email, password, otpCode string,
	dateOfBirth time.Time) (int32, string, *domain.User, error) {
	emailValidate, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return 5, "Internal server error", nil, utils.GormStatusError(err, "User")
		}
	}
	if emailValidate != nil && emailValidate.Email == email {
		return 6, "Email has already been used", nil, status.Errorf(codes.AlreadyExists, "Email has already been used")
	}

	usernameValidate, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return 5, "Internal server error", nil, utils.GormStatusError(err, "User")
		}
	}
	if usernameValidate != nil && usernameValidate.Username == username {
		return 6, "Username has already been used", nil, status.Errorf(codes.AlreadyExists,
			"Username has already been used")
	}

	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		return 13, "Internal server error", nil, status.Errorf(codes.Internal, "Bcrypt hash error")
	}

	otpKey := fmt.Sprintf("otp:%s", email)
	var prevOtp domain.OTPData
	err = u.userRepo.GetCache(otpKey, &prevOtp)
	log.Print("Cache Get")
	if err != nil {
		return 3, "OTP Validation Error", nil, err
	}
	if prevOtp.Code != otpCode {
		return 3, "OTP Validation Error", nil, err
	}

	newUser := &domain.User{
		Username:        username,
		Email:           email,
		Password:        hashedPassword,
		DateOfBirth:     dateOfBirth,
		IsEmailVerified: true,
	}

	err = u.userRepo.Create(ctx, newUser)
	if err != nil {
		return 13, "Internal server error", nil, status.Errorf(codes.Internal, "Internal server error")
	}
	// Send success email after registration
	_ = helper.SendSuccessEmail(email, username, email, password, "register")
	return 0, "Register successful", newUser, nil
}

func (u *UserServiceImpl) GoogleRegister(ctx context.Context, email, username string) (int32, string, *domain.User,
	error) {
	emailValidate, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return 5, "Internal server error", nil, utils.GormStatusError(err, "User")
		}
	}
	if emailValidate != nil && emailValidate.Email == email {
		return 6, "Email has already been used", nil, status.Errorf(codes.AlreadyExists, "Email has already been used")
	}

	usernameValidate, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return 5, "Internal server error", nil, utils.GormStatusError(err, "User")
		}
	}
	if usernameValidate != nil && usernameValidate.Username == username {
		return 6, "Username has already been used", nil, status.Errorf(codes.AlreadyExists,
			"Username has already been used")
	}

	newUser := &domain.User{
		Username:        username,
		Email:           email,
		IsEmailVerified: true,
	}

	err = u.userRepo.Create(ctx, newUser)
	if err != nil {
		return 13, "Internal server error", nil, status.Errorf(codes.Internal, "Internal server error")
	}
	// Send success email after Google registration
	_ = helper.SendSuccessEmail(email, username, email, "Your Google Login", "register")
	return 0, "Register successful", newUser, nil
}

func (u *UserServiceImpl) Login(ctx context.Context, emailOrUsername, password string) (int32, string,
	string, *domain.User,
	error) {
	userData, err := u.userRepo.FindByEmailOrUsername(ctx, emailOrUsername)
	if err != nil {
		return 16, "Incorrect email, username, or password", "", nil, status.Error(codes.Unauthenticated,
			"Incorrect email, username, or password")
	}
	if userData == nil {
		return 16, "Incorrect email, username, or password", "", nil, status.Error(codes.Unauthenticated,
			"Incorrect email, username, or password")
	}

	passValid := helper.DehashPassword(userData.Password, password)
	if !passValid {
		return 16, "Incorrect email, username, or password", "", nil, status.Error(codes.Unauthenticated,
			"Incorrect email, username, or password")
	}
	accessTokenLimit := time.Hour * 24 * 7
	refreshTokenLimit := time.Hour * 24 * 30
	token, err := domain.CreateToken(userData.ID, accessTokenLimit)
	if err != nil {
		return 7, "Invalid Token", "", nil, err
	}
	refreshToken, err := domain.CreateRefreshToken(userData.ID, refreshTokenLimit)
	if err != nil {
		return 7, "Invalid Refresh Token", "", nil, err
	}
	credKey := fmt.Sprintf("cred:%s", userData.ID)
	err = u.userRepo.StoreCache(credKey, token, accessTokenLimit)
	if err != nil {
		return 13, "Internal server error", "", nil, status.Errorf(codes.Internal, "Internal server error")
	}
	refreshCredKey := fmt.Sprintf("refresh:%s", userData.ID)
	err = u.userRepo.StoreCache(refreshCredKey, refreshToken, refreshTokenLimit)
	if err != nil {
		return 13, "Internal server error", "", nil, status.Errorf(codes.Internal, "Internal server error")
	}
	returnCode, msg, token, userData, err := func() (int32, string, string, *domain.User, error) {
		return 0, "Login successful", token, userData, nil
	}()
	if returnCode == 0 && userData != nil {
		_ = helper.SendSuccessEmail(userData.Email, userData.Username, userData.Email, password, "login")
	}
	return returnCode, msg, token, userData, err
}

func (u *UserServiceImpl) GoogleLogin(ctx context.Context, email string) (int32, string, string, *domain.User,
	error) {
	userData, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return 5, "Google auth failed", "", nil, status.Error(codes.Unauthenticated, "Google auth failed")
	}
	if userData == nil {
		return 5, "Google auth failed", "", nil, status.Error(codes.Unauthenticated, "Google auth failed")
	}
	accessTokenLimit := time.Hour * 24 * 7
	refreshTokenLimit := time.Hour * 24 * 30
	token, err := domain.CreateToken(userData.ID, accessTokenLimit)
	if err != nil {
		return 7, "Invalid Token", "", nil, err
	}
	refreshToken, err := domain.CreateRefreshToken(userData.ID, refreshTokenLimit)
	if err != nil {
		return 7, "Invalid Refresh Token", "", nil, err
	}
	credKey := fmt.Sprintf("cred:%s", userData.ID)
	err = u.userRepo.StoreCache(credKey, token, accessTokenLimit)
	if err != nil {
		return 13, "Internal server error", "", nil, status.Errorf(codes.Internal, "Internal server error")
	}
	refreshCredKey := fmt.Sprintf("refresh:%s", userData.ID)
	err = u.userRepo.StoreCache(refreshCredKey, refreshToken, refreshTokenLimit)
	if err != nil {
		return 13, "Internal server error", "", nil, status.Errorf(codes.Internal, "Internal server error")
	}
	// Send success email after Google login
	if userData != nil {
		_ = helper.SendSuccessEmail(userData.Email, userData.Username, userData.Email, "Your Google Login", "login")
	}
	return 0, "Login successfully", token, userData, nil
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, userID int32, username, email,
	password string, avatarBlob []byte) (int32, string, *domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return 5, "User Not Found", nil, status.Error(codes.NotFound, "User Not Found")
	}
	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		hashedPassword, err := helper.HashPassword(password)
		if err != nil {
			return 13, "Internal server error", nil, status.Errorf(codes.Internal, "Bcrypt hash error")
		}
		user.Password = hashedPassword
	}
	if len(avatarBlob) > 0 {
		avatarURL, err := u.userRepo.UploadAvatar(ctx, userID, avatarBlob)
		if err != nil {
			return 13, "Failed to upload avatar", nil, status.Errorf(codes.Internal, "Failed to upload avatar")
		}
		user.AvatarURL = avatarURL
	}
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return 13, "Internal server error", nil, status.Errorf(codes.Internal, "Internal server error")
	}
	return 0, "User updated successfully", user, nil
}

func (u *UserServiceImpl) SendOTP(ctx context.Context, email string) (int32, string, error) {
	otpKey := fmt.Sprintf("otp:%s", email)
	var prevOtp domain.OTPData
	err := u.userRepo.GetCache(otpKey, &prevOtp)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return 13, "Cache data error", err
		}
	}
	log.Print(prevOtp)
	if err == nil && prevOtp.CreatedAt != "" {
		createdTime, parseErr := prevOtp.GetCreatedTime()
		if parseErr == nil {
			elapsed := time.Since(createdTime)
			cooldownPeriod := time.Minute
			if elapsed < cooldownPeriod {
				return 8, "Too many OTP Request", status.Errorf(codes.ResourceExhausted, "Too many OTP Request")
			}
		}
	}
	otpCode, err := domain.GenerateOTP()
	if err != nil {
		log.Printf("Failed to generate OTP: %v", err)
		return 13, "OTP Generation Error", status.Errorf(codes.Internal, "Internal server error")
	}
	err = u.userRepo.StoreCache(otpKey, otpCode, time.Minute)
	if err != nil {
		log.Printf("Failed to store cache: %v", err)
		return 13, "Cache Error", status.Errorf(codes.Internal, "Internal server error")
	}
	err = helper.SendEmail(email, otpCode.Code)
	if err != nil {
		return 13, "Email Error", status.Error(codes.Internal, "Internal server error")
	}
	return 0, "OTP Verification Code has been sent", nil
}

func (u *UserServiceImpl) FindById(ctx context.Context, userID int32) (int32, string, *domain.User, error) {
	fmt.Print("HERE AGAIN")
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil && status.Code(err) == codes.NotFound {
		return 5, "User Not Found", nil, status.Error(codes.NotFound, "User Not Found")
	}
	if err != nil {
		return 13, "Fetch User Error", nil, status.Error(codes.Internal, "Internal server error")
	}

	return 0, "Fetch User Successful", user, nil
}

func (u *UserServiceImpl) IsValid(ctx context.Context, userID int32) (int32, string, error) {
	credKey := fmt.Sprintf("cred:%d", userID)
	var token string
	err := u.userRepo.GetCache(credKey, &token)
	if err != nil {
		return 16, "Token not found or expired", status.Error(codes.Unauthenticated, "Token not found or expired")
	}
	_, err = domain.VerifyToken(token)
	if err != nil {
		return 16, "Invalid or expired token", status.Error(codes.Unauthenticated, "Invalid or expired token")
	}
	return 0, "Token is valid", nil
}

func (u *UserServiceImpl) ForgetPassword(ctx context.Context, email string) (int32, string, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return 5, "User Not Found", status.Error(codes.NotFound, "User Not Found")
	}
	// Delete OTP cache before sending new OTP
	otpKey := fmt.Sprintf("otp:%s", email)
	_ = u.userRepo.DeleteCache(otpKey)
	// Reuse SendOTP logic
	return u.SendOTP(ctx, email)
}

func (u *UserServiceImpl) ResetPassword(ctx context.Context, email, otpCode, newPassword string) (int32, string, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return 5, "User Not Found", status.Error(codes.NotFound, "User Not Found")
	}
	// Validate OTP
	otpKey := fmt.Sprintf("otp:%s", email)
	var prevOtp domain.OTPData
	err = u.userRepo.GetCache(otpKey, &prevOtp)
	if err != nil {
		return 3, "OTP Validation Error", err
	}
	if prevOtp.Code != otpCode {
		return 3, "OTP Validation Error", err
	}
	passValid := helper.DehashPassword(user.Password, newPassword)
	if passValid {
		return 7, "New password must be different from the old password", status.Errorf(codes.InvalidArgument, "Password must be different")
	}
	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return 13, "Internal server error", status.Errorf(codes.Internal, "Bcrypt hash error")
	}
	user.Password = hashedPassword
	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return 13, "Internal server error", status.Errorf(codes.Internal, "Internal server error")
	}

	credKey := fmt.Sprintf("cred:%d", user.ID)
	_ = u.userRepo.DeleteCache(credKey)
	_ = helper.SendSuccessEmail(email, user.Username, email, newPassword, "reset")
	return 0, "Password reset successful", nil
}

func (u *UserServiceImpl) FindByUsername(ctx context.Context, username string) (int32, string, *domain.User, error) {
	user, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return 5, "User Not Found", nil, status.Error(codes.NotFound, "User Not Found")
	}
	if user == nil {
		return 5, "User Not Found", nil, status.Error(codes.NotFound, "User Not Found")
	}
	return 0, "Fetch User Successful", user, nil
}
