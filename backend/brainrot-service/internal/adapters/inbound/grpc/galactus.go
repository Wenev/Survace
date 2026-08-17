package grpc

import (
	"context"
	"fmt"
	"github.com/Wenev/Survace/proto/gen/controller"
	"github.com/Wenev/Survace/proto/gen/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GalactusClient struct {
	Client controller.GalactusControllerClient
	Social controller.SocialServiceClient
}

func NewGalactusClient(address string) (*GalactusClient, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Galactus: %v", err)
	}
	client := controller.NewGalactusControllerClient(conn)
	return &GalactusClient{Client: client}, nil
}

func (g *GalactusClient) FindByUserId(ctx context.Context, userId int32) (*dto.FindByUserIdResponse, error) {
	resp, err := g.Client.FindByUserId(ctx, &dto.FindByUserIdRequest{Id: userId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) ListFriend(ctx context.Context, userId int32) (*dto.ListFriendsResponse, error) {
	resp, err := g.Social.ListFriends(ctx, &dto.ListFriendsRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (g *GalactusClient) ListFollowing(ctx context.Context, userId int32) (*dto.ListFollowingResponse, error) {
	resp, err := g.Social.ListFollowing(ctx, &dto.ListFollowingRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalactusClient) IsValid(ctx context.Context, userId int32) (*dto.IsValidResponse, error) {
	resp, err := g.Client.IsValid(ctx, &dto.IsValidRequest{UserId: userId})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
