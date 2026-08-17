package authentication

import (
	"context"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

type GoogleUserInfo struct {
	Email    string
	Name     string
	Picture  string
	GoogleID string
}

func FetchGoogleOAuth(ctx context.Context, googleToken string) (*GoogleUserInfo, error) {
	clientID := os.Getenv("GOOGLE_OATH_CLIENT_ID")
	if clientID == "" {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	payload, err := idtoken.Validate(ctx, googleToken, clientID)
	if err != nil {
		return nil, err
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)
	googleID, _ := payload.Claims["sub"].(string)

	if email == "" {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &GoogleUserInfo{
		Email:    email,
		Name:     name,
		Picture:  picture,
		GoogleID: googleID,
	}, nil
}
