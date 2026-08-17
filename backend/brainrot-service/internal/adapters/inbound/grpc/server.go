package grpc

import (
	"context"
	"fmt"
	out "github.com/Wenev/Survace/brainrot-service/ports/in"
	"github.com/Wenev/Survace/proto/gen/controller"
	"github.com/Wenev/Survace/proto/gen/dto"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
	"time"
)

type GrpcServer struct {
	grpcPort           int
	server             *grpc.Server
	videoService       out.VideoService
	playlistService    out.PlaylistService // <-- Add this line
	commentService     out.CommentService
	replyService       out.ReplyService
	likeService        out.LikeService
	likeCommentService out.LikeCommentService
	likeReplyService   out.LikeReplyService
	controller.BrainrotServiceServer
	GalactusClient *GalactusClient
}

func (g *GrpcServer) middleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	fmt.Print("HII")
	methodsRequiringUser := map[string]bool{
		"/controller.BrainrotService/UploadVideo":     true,
		"/controller.BrainrotService/AddComment":      true,
		"/controller.BrainrotService/DeleteUserVideo": true,
		"/controller.BrainrotService/DeleteComment":   true,
		"/controller.BrainrotService/AddReply":        true,
		"/controller.BrainrotService/DeleteReply":     true,
		"/controller.BrainrotService/LikeVideo":       true,
		"/controller.BrainrotService/UnlikeVideo":     true,
		"/controller.BrainrotService/LikeComment":     true,
		"/controller.BrainrotService/UnlikeComment":   true,
		"/controller.BrainrotService/LikeReply":       true,
		"/controller.BrainrotService/UnlikeReply":     true,
		"/controller.BrainrotService/WatchVideo":      true,
		"/controller.BrainrotService/GetRandomFeed":   true,
	}

	if methodsRequiringUser[info.FullMethod] {
		var userId int32
		switch r := req.(type) {
		case *dto.UploadVideoRequest:
			userId = r.UserId
		case *dto.AddCommentRequest:
			userId = r.UserId
		case *dto.DeleteUserVideoRequest:
			userId = r.UserId
		case *dto.AddReplyRequest:
			userId = r.UserId
		case *dto.LikeVideoRequest:
			userId = r.UserId
		case *dto.UnlikeVideoRequest:
			userId = r.UserId
		case *dto.LikeCommentRequest:
			userId = r.UserId
		case *dto.UnlikeCommentRequest:
			userId = r.UserId
		case *dto.LikeReplyRequest:
			userId = r.UserId
		case *dto.UnlikeReplyRequest:
			userId = r.UserId
		case *dto.WatchVideoRequest:
			userId = r.UserId
		case *dto.GetRandomFeedRequest:
			userId = r.UserId
		}

		if userId == 0 {
			return nil, fmt.Errorf("user_id is required for this operation")
		}

		resp, err := g.GalactusClient.IsValid(ctx, userId)
		if err != nil || resp == nil || resp.StatusCode != 0 {
			return nil, fmt.Errorf("user not found or invalid")
		}
	}

	return handler(ctx, req)
}

func NewGrpcServer(grpcPort int, videoService out.VideoService, playlistService out.PlaylistService, commentService out.CommentService,
	replyService out.ReplyService, likeService out.LikeService, likeCommentService out.LikeCommentService, likeReplyService out.LikeReplyService, galactusClient *GalactusClient) *GrpcServer {
	gs := &GrpcServer{
		grpcPort:           grpcPort,
		videoService:       videoService,
		playlistService:    playlistService, // <-- Add this line
		commentService:     commentService,
		replyService:       replyService,
		likeService:        likeService,
		likeCommentService: likeCommentService,
		likeReplyService:   likeReplyService,
		GalactusClient:     galactusClient,
	}
	gs.server = grpc.NewServer(
		grpc.MaxRecvMsgSize(100*1024*1024),
		grpc.UnaryInterceptor(gs.middleware),
	)
	return gs
}

func (g *GrpcServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", g.grpcPort))
	if err != nil {
		fmt.Print("HELP")
		return err
	}
	controller.RegisterBrainrotServiceServer(g.server, g)
	err = g.server.Serve(listener)
	if err != nil {
		return err
	}
	return nil
}

func (g *GrpcServer) Stop() {
	g.server.GracefulStop()
}

func (g *GrpcServer) UploadVideo(ctx context.Context,
	in *dto.UploadVideoRequest) (*dto.UploadVideoResponse, error) {
	var postDate *time.Time
	if in.PostDate != nil && in.PostDate.Seconds != 0 {
		t := in.PostDate.AsTime()
		postDate = &t
	}
	code, message, video, err := g.videoService.UploadVideo(ctx, in.Title, in.File, in.UserId, in.EnableComment, in.Visibility, in.Thumbnail, in.IsDraft, postDate)
	if err != nil {
		return nil, err
	}

	var protoVideo *dto.Video
	if video != nil {
		protoVideo = &dto.Video{
			Id:            video.ID,
			ObjectName:    video.ObjectName,
			Url:           video.URL,
			Title:         video.Title,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			CreatedAt:     timestamppb.New(video.CreatedAt),
			ThumbnailUrl:  video.ThumbnailUrl,
			PostedAt:      timestamppb.New(video.PostedAt),
			IsDraft:       video.IsDraft,
		}
	}

	resp := &dto.UploadVideoResponse{
		Code:    code,
		Message: message,
		Video:   protoVideo,
	}
	return resp, nil
}

func (g *GrpcServer) GetVideoByID(ctx context.Context, in *dto.GetVideoByIDRequest) (*dto.GetVideoByIDResponse, error) {
	code, message, video, likeCount, commentCount, viewCount, err := g.videoService.GetVideoByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	vid := &dto.VideoFeed{
		Id:            video.ID,
		UserId:        video.UserID,
		ObjectName:    video.ObjectName,
		Url:           video.URL,
		Title:         video.Title,
		EnableComment: video.EnableComment,
		Visibility:    video.Visibility,
		CreatedAt:     timestamppb.New(video.CreatedAt),
		LikeCount:     int32(likeCount),
		CommentCount:  int32(commentCount),
		ViewCount:     int32(viewCount),
	}

	resp := &dto.GetVideoByIDResponse{
		Code:    code,
		Message: message,
		Video:   vid,
	}
	return resp, nil
}

func (g *GrpcServer) GetUserVideo(ctx context.Context, in *dto.GetUserVideoRequest) (*dto.GetUserVideoResponse, error) {
	code, message, videos, likeCounts, commentCounts, viewCounts, err := g.videoService.GetUserVideo(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	protoVideos := make([]*dto.VideoFeed, 0, len(videos))
	for i, video := range videos {
		vid := &dto.VideoFeed{
			Id:            video.ID,
			UserId:        video.UserID,
			ObjectName:    video.ObjectName,
			Url:           video.URL,
			Title:         video.Title,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			CreatedAt:     timestamppb.New(video.CreatedAt),
			LikeCount:     int32(likeCounts[i]),
			CommentCount:  int32(commentCounts[i]),
			ViewCount:     int32(viewCounts[i]),
		}
		protoVideos = append(protoVideos, vid)
	}
	fmt.Print("HFFHFHfdfsfdsssfsdf")
	resp := &dto.GetUserVideoResponse{
		Code:    code,
		Message: message,
		Videos:  protoVideos,
	}
	return resp, nil
}

func (g *GrpcServer) GetUserVideoAndDraft(ctx context.Context, in *dto.GetUserVideoAndDraftRequest) (*dto.GetUserVideoAndDraftResponse, error) {
	code, message, videos, likeCounts, commentCounts, viewCounts, err := g.videoService.GetUserVideoAndDraft(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	protoVideos := make([]*dto.VideoFeed, 0, len(videos))
	for i, video := range videos {
		vid := &dto.VideoFeed{
			Id:            video.ID,
			UserId:        video.UserID,
			ObjectName:    video.ObjectName,
			Url:           video.URL,
			Title:         video.Title,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			CreatedAt:     timestamppb.New(video.CreatedAt),
			LikeCount:     int32(likeCounts[i]),
			CommentCount:  int32(commentCounts[i]),
			ViewCount:     int32(viewCounts[i]),
			ThumbnailUrl:  video.ThumbnailUrl,
			PostedAt:      timestamppb.New(video.PostedAt),
			IsDraft:       video.IsDraft,
		}
		protoVideos = append(protoVideos, vid)
	}
	return &dto.GetUserVideoAndDraftResponse{
		Code:    code,
		Message: message,
		Videos:  protoVideos,
	}, nil
}

func (g *GrpcServer) DeleteUserVideo(ctx context.Context, in *dto.DeleteUserVideoRequest) (*dto.DeleteUserVideoResponse, error) {
	code, message, err := g.videoService.DeleteUserVideo(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	resp := &dto.DeleteUserVideoResponse{
		Code:    code,
		Message: message,
	}
	return resp, nil
}

func (g *GrpcServer) AddComment(ctx context.Context, in *dto.AddCommentRequest) (*dto.AddCommentResponse, error) {
	code, message, err := g.commentService.AddComment(ctx, in.UserId, in.VideoId, in.Text)
	return &dto.AddCommentResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) GetCommentsByVideo(ctx context.Context, in *dto.GetCommentsByVideoRequest) (*dto.GetCommentsByVideoResponse, error) {
	code, message, comments, replyCounts, likeCounts, err := g.commentService.GetCommentsByVideo(ctx, in.VideoId)
	if err != nil {
		return nil, err
	}
	protoComments := make([]*dto.CommentFeed, 0, len(comments))
	for i, c := range comments {
		protoComments = append(protoComments, &dto.CommentFeed{
			Id:         c.ID,
			Content:    c.Content,
			VideoId:    c.VideoID,
			UserId:     c.UserID,
			LikeCount:  int32(likeCounts[i]),
			ReplyCount: int32(replyCounts[i]),
		})
	}
	return &dto.GetCommentsByVideoResponse{
		Code:     code,
		Message:  message,
		Comments: protoComments,
	}, nil
}

func (g *GrpcServer) DeleteComment(ctx context.Context, in *dto.DeleteCommentRequest) (*dto.DeleteCommentResponse, error) {
	code, message, err := g.commentService.DeleteComment(ctx, in.CommentId)
	return &dto.DeleteCommentResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) AddReply(ctx context.Context, in *dto.AddReplyRequest) (*dto.AddReplyResponse, error) {
	code, message, err := g.replyService.AddReply(ctx, in.UserId, in.CommentId, in.Text)
	return &dto.AddReplyResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) GetRepliesByComment(ctx context.Context, in *dto.GetRepliesByCommentRequest) (*dto.GetRepliesByCommentResponse, error) {
	code, message, replies, likeCounts, err := g.replyService.GetRepliesByComment(ctx, in.CommentId)
	if err != nil {
		return nil, err
	}
	protoReplies := make([]*dto.ReplyFeed, 0, len(replies))
	for i, r := range replies {
		protoReplies = append(protoReplies, &dto.ReplyFeed{
			Id:        r.ID,
			Content:   r.Content,
			CommentId: r.CommentID,
			UserId:    r.UserID,
			LikeCount: int32(likeCounts[i]),
		})
	}
	return &dto.GetRepliesByCommentResponse{
		Code:    code,
		Message: message,
		Replies: protoReplies,
	}, nil
}

func (g *GrpcServer) DeleteReply(ctx context.Context, in *dto.DeleteReplyRequest) (*dto.DeleteReplyResponse, error) {
	code, message, err := g.replyService.DeleteReply(ctx, in.ReplyId)
	return &dto.DeleteReplyResponse{
		Code:    code,
		Message: message,
	}, err
}

// --- Like Video Endpoints ---
func (g *GrpcServer) LikeVideo(ctx context.Context, in *dto.LikeVideoRequest) (*dto.LikeVideoResponse, error) {
	code, message, err := g.likeService.LikeVideo(ctx, in.UserId, in.VideoId)
	return &dto.LikeVideoResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) UnlikeVideo(ctx context.Context, in *dto.UnlikeVideoRequest) (*dto.UnlikeVideoResponse, error) {
	code, message, err := g.likeService.UnlikeVideo(ctx, in.UserId, in.VideoId)
	return &dto.UnlikeVideoResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) IsVideoLiked(ctx context.Context, in *dto.IsVideoLikedRequest) (*dto.IsVideoLikedResponse, error) {
	liked, err := g.likeService.IsVideoLiked(ctx, in.UserId, in.VideoId)
	return &dto.IsVideoLikedResponse{
		Liked: liked,
	}, err
}

func (g *GrpcServer) LikeComment(ctx context.Context, in *dto.LikeCommentRequest) (*dto.LikeCommentResponse, error) {
	code, message, err := g.likeCommentService.LikeComment(ctx, in.UserId, in.CommentId)
	return &dto.LikeCommentResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) UnlikeComment(ctx context.Context, in *dto.UnlikeCommentRequest) (*dto.UnlikeCommentResponse, error) {
	code, message, err := g.likeCommentService.UnlikeComment(ctx, in.UserId, in.CommentId)
	return &dto.UnlikeCommentResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) IsCommentLiked(ctx context.Context, in *dto.IsCommentLikedRequest) (*dto.IsCommentLikedResponse, error) {
	liked, err := g.likeCommentService.IsCommentLiked(ctx, in.UserId, in.CommentId)
	return &dto.IsCommentLikedResponse{
		Liked: liked,
	}, err
}

func (g *GrpcServer) LikeReply(ctx context.Context, in *dto.LikeReplyRequest) (*dto.LikeReplyResponse, error) {
	code, message, err := g.likeReplyService.LikeReply(ctx, in.ReplyId, in.UserId)
	return &dto.LikeReplyResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) UnlikeReply(ctx context.Context, in *dto.UnlikeReplyRequest) (*dto.UnlikeReplyResponse, error) {
	code, message, err := g.likeReplyService.UnlikeReply(ctx, in.ReplyId, in.UserId)
	return &dto.UnlikeReplyResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) IsReplyLiked(ctx context.Context, in *dto.IsReplyLikedRequest) (*dto.IsReplyLikedResponse, error) {
	liked, err := g.likeReplyService.IsReplyLiked(ctx, in.ReplyId, in.UserId)
	return &dto.IsReplyLikedResponse{
		Liked: liked,
	}, err
}

func (g *GrpcServer) WatchVideo(ctx context.Context, in *dto.WatchVideoRequest) (*dto.WatchVideoResponse, error) {
	code, message, err := g.videoService.WatchVideo(ctx, in.UserId, in.VideoId)
	return &dto.WatchVideoResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) GetRandomFeed(ctx context.Context, in *dto.GetRandomFeedRequest) (*dto.GetRandomFeedResponse, error) {
	var userIdPtr *int32
	if in.UserId != 0 {
		userIdPtr = &in.UserId
	}
	code, message, videos, likeCounts, commentCounts, viewCounts, err := g.videoService.GetRandomFeed(ctx, userIdPtr, int(in.Limit), int(in.Offset))
	if err != nil {
		return &dto.GetRandomFeedResponse{
			Code:    code,
			Message: message,
			Videos:  nil,
		}, err
	}
	if len(videos) == 0 {
		return &dto.GetRandomFeedResponse{
			Code:    code,
			Message: message,
			Videos:  nil,
		}, nil
	}
	protoVideos := make([]*dto.VideoFeed, 0, len(videos))
	for i, video := range videos {
		protoVideos = append(protoVideos, &dto.VideoFeed{
			Id:            video.ID,
			UserId:        video.UserID,
			ObjectName:    video.ObjectName,
			Url:           video.URL,
			Title:         video.Title,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			CreatedAt:     timestamppb.New(video.CreatedAt),
			LikeCount:     int32(likeCounts[i]),
			CommentCount:  int32(commentCounts[i]),
			ViewCount:     int32(viewCounts[i]),
			ThumbnailUrl:  video.ThumbnailUrl, // Ensure thumbnail is set
			PostedAt:      timestamppb.New(video.PostedAt),
			IsDraft:       video.IsDraft,
		})
	}
	return &dto.GetRandomFeedResponse{
		Code:    code,
		Message: message,
		Videos:  protoVideos,
	}, nil
}

func (g *GrpcServer) GetRandomFeedLoggedOut(ctx context.Context,
	in *dto.GetRandomFeedLoggedOutRequest) (*dto.GetRandomFeedLoggedOutResponse, error) {
	code, message, videos, likeCounts, commentCounts, viewCounts, err := g.videoService.GetRandomFeedLoggedOut(ctx, int(in.Limit), int(in.Offset))
	if err != nil {
		return &dto.GetRandomFeedLoggedOutResponse{
			Code:    code,
			Message: message,
			Videos:  nil,
		}, err
	}
	if len(videos) == 0 {
		return &dto.GetRandomFeedLoggedOutResponse{
			Code:    code,
			Message: message,
			Videos:  nil,
		}, nil
	}
	protoVideos := make([]*dto.VideoFeed, 0, len(videos))
	for i, video := range videos {
		protoVideos = append(protoVideos, &dto.VideoFeed{
			Id:            video.ID,
			UserId:        video.UserID,
			ObjectName:    video.ObjectName,
			Url:           video.URL,
			Title:         video.Title,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			CreatedAt:     timestamppb.New(video.CreatedAt),
			LikeCount:     int32(likeCounts[i]),
			CommentCount:  int32(commentCounts[i]),
			ViewCount:     int32(viewCounts[i]),
			ThumbnailUrl:  video.ThumbnailUrl, // Ensure thumbnail is set
			PostedAt:      timestamppb.New(video.PostedAt),
			IsDraft:       video.IsDraft,
		})
	}
	return &dto.GetRandomFeedLoggedOutResponse{
		Code:    code,
		Message: message,
		Videos:  protoVideos,
	}, nil
}

func (g *GrpcServer) SearchVideo(ctx context.Context, in *dto.SearchVideoRequest) (*dto.SearchVideoResponse, error) {
	code, message, videos, likeCounts, commentCounts, viewCounts, err := g.videoService.SearchVideo(ctx, in.Query, in.Threshold, int(in.Limit), int(in.Offset))
	if err != nil {
		return &dto.SearchVideoResponse{
			Code:    code,
			Message: message,
			Videos:  nil,
		}, err
	}
	protoVideos := make([]*dto.VideoFeed, 0, len(videos))
	for i, video := range videos {
		protoVideos = append(protoVideos, &dto.VideoFeed{
			Id:            video.ID,
			UserId:        video.UserID,
			ObjectName:    video.ObjectName,
			Url:           video.URL,
			Title:         video.Title,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			CreatedAt:     timestamppb.New(video.CreatedAt),
			LikeCount:     int32(likeCounts[i]),
			CommentCount:  int32(commentCounts[i]),
			ViewCount:     int32(viewCounts[i]),
		})
	}
	return &dto.SearchVideoResponse{
		Code:    code,
		Message: message,
		Videos:  protoVideos,
	}, nil
}

func (g *GrpcServer) DeleteVideoByID(ctx context.Context, in *dto.DeleteVideoByIDRequest) (*dto.DeleteVideoByIDResponse, error) {
	code, message, err := g.videoService.DeleteVideoByID(ctx, in.VideoId)
	if err != nil {
		return &dto.DeleteVideoByIDResponse{
			Code:    code,
			Message: message,
		}, err
	}
	return &dto.DeleteVideoByIDResponse{
		Code:    code,
		Message: message,
	}, nil
}

func (g *GrpcServer) UpdateVideo(ctx context.Context, in *dto.UpdateVideoRequest) (*dto.UpdateVideoResponse, error) {
	var postedAt *time.Time
	if in.PostedAt != nil && in.PostedAt.Seconds != 0 {
		t := in.PostedAt.AsTime()
		postedAt = &t
	}
	code, message, video, err := g.videoService.UpdateVideo(
		ctx,
		in.Id,
		in.UserId,
		in.Title,
		in.Url,
		in.ObjectName,
		in.EnableComment,
		in.Visibility,
		in.Thumbnail, // Pass thumbnail bytes
		postedAt,
		in.IsDraft,
	)
	if err != nil {
		return &dto.UpdateVideoResponse{
			Code:    code,
			Message: message,
			Video:   nil,
		}, err
	}
	var protoVideo *dto.Video
	if video != nil {
		protoVideo = &dto.Video{
			Id:            video.ID,
			UserId:        video.UserID,
			Title:         video.Title,
			Url:           video.URL,
			ObjectName:    video.ObjectName,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			ThumbnailUrl:  video.ThumbnailUrl,
			PostedAt:      timestamppb.New(video.PostedAt),
			IsDraft:       video.IsDraft,
			CreatedAt:     timestamppb.New(video.CreatedAt),
		}
	}
	return &dto.UpdateVideoResponse{
		Code:    code,
		Message: message,
		Video:   protoVideo,
	}, nil
}

func (g *GrpcServer) FriendVideo(ctx context.Context, in *dto.FriendVideoRequest) (*dto.FriendVideoResponse,
	error) {
	code, message, videos, likeCounts, commentCounts, viewCounts, err := g.videoService.FriendVideo(ctx, in.UserId,
		int(in.Limit), int(in.Offset))
	if err != nil {
		return &dto.FriendVideoResponse{
			Code:    code,
			Message: message,
			Videos:  nil,
		}, err
	}
	protoVideos := make([]*dto.VideoFeed, 0, len(videos))
	for i, video := range videos {
		protoVideos = append(protoVideos, &dto.VideoFeed{
			Id:            video.ID,
			UserId:        video.UserID,
			ObjectName:    video.ObjectName,
			Url:           video.URL,
			Title:         video.Title,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			CreatedAt:     timestamppb.New(video.CreatedAt),
			LikeCount:     int32(likeCounts[i]),
			CommentCount:  int32(commentCounts[i]),
			ViewCount:     int32(viewCounts[i]),
			ThumbnailUrl:  video.ThumbnailUrl,
			PostedAt:      timestamppb.New(video.PostedAt),
			IsDraft:       video.IsDraft,
		})
	}
	return &dto.FriendVideoResponse{
		Code:    code,
		Message: message,
		Videos:  protoVideos,
	}, nil
}

func (g *GrpcServer) FollowingVideo(ctx context.Context, in *dto.FollowingVideoRequest) (*dto.FollowingVideoResponse, error) {
	code, message, videos, likeCounts, commentCounts, viewCounts, err := g.videoService.FollowingVideo(ctx, in.UserId, int(in.Limit), int(in.Offset))
	if err != nil {
		return &dto.FollowingVideoResponse{
			Code:    code,
			Message: message,
			Videos:  nil,
		}, err
	}
	protoVideos := make([]*dto.VideoFeed, 0, len(videos))
	for i, video := range videos {
		protoVideos = append(protoVideos, &dto.VideoFeed{
			Id:            video.ID,
			UserId:        video.UserID,
			ObjectName:    video.ObjectName,
			Url:           video.URL,
			Title:         video.Title,
			EnableComment: video.EnableComment,
			Visibility:    video.Visibility,
			CreatedAt:     timestamppb.New(video.CreatedAt),
			LikeCount:     int32(likeCounts[i]),
			CommentCount:  int32(commentCounts[i]),
			ViewCount:     int32(viewCounts[i]),
			ThumbnailUrl:  video.ThumbnailUrl,
			PostedAt:      timestamppb.New(video.PostedAt),
			IsDraft:       video.IsDraft,
		})
	}
	return &dto.FollowingVideoResponse{
		Code:    code,
		Message: message,
		Videos:  protoVideos,
	}, nil
}

func (g *GrpcServer) CreatePlaylist(ctx context.Context, in *dto.CreatePlaylistRequest) (*dto.CreatePlaylistResponse, error) {
	code, message, playlist, err := g.playlistService.CreatePlaylist(ctx, in.UserId, in.Title)
	if err != nil {
		return &dto.CreatePlaylistResponse{
			Code:     code,
			Message:  message,
			Playlist: nil,
		}, err
	}
	var protoPlaylist *dto.Playlist
	if playlist != nil {
		protoPlaylist = &dto.Playlist{
			Id:     playlist.ID,
			UserId: playlist.UserID,
			Title:  playlist.Title,
			Videos: nil, // Videos can be filled if needed
		}
	}
	return &dto.CreatePlaylistResponse{
		Code:     code,
		Message:  message,
		Playlist: protoPlaylist,
	}, nil
}

// --- Playlist Endpoints ---

func (g *GrpcServer) DeletePlaylist(ctx context.Context, in *dto.DeletePlaylistRequest) (*dto.DeletePlaylistResponse, error) {
	code, message, err := g.playlistService.DeletePlaylist(ctx, in.PlaylistId, in.UserId)
	return &dto.DeletePlaylistResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) AddVideoToPlaylist(ctx context.Context, in *dto.AddVideoToPlaylistRequest) (*dto.AddVideoToPlaylistResponse, error) {
	code, message, err := g.playlistService.AddVideoToPlaylist(ctx, in.PlaylistId, in.VideoId, in.Order)
	return &dto.AddVideoToPlaylistResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) RemoveVideoFromPlaylist(ctx context.Context, in *dto.RemoveVideoFromPlaylistRequest) (*dto.RemoveVideoFromPlaylistResponse, error) {
	code, message, err := g.playlistService.RemoveVideoFromPlaylist(ctx, in.PlaylistId, in.VideoId)
	return &dto.RemoveVideoFromPlaylistResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) ReorderVideo(ctx context.Context, in *dto.ReorderVideoRequest) (*dto.ReorderVideoResponse, error) {
	code, message, err := g.playlistService.ReorderVideo(ctx, in.PlaylistId, in.VideoId, in.NewOrder)
	return &dto.ReorderVideoResponse{
		Code:    code,
		Message: message,
	}, err
}

func (g *GrpcServer) GetUserPlaylists(ctx context.Context, in *dto.GetUserPlaylistsRequest) (*dto.GetUserPlaylistsResponse, error) {
	code, message, playlists, err := g.playlistService.GetUserPlaylists(ctx, in.UserId)
	if err != nil {
		return &dto.GetUserPlaylistsResponse{
			Code:      code,
			Message:   message,
			Playlists: nil,
		}, err
	}
	protoPlaylists := make([]*dto.Playlist, 0, len(playlists))
	for _, playlist := range playlists {
		protoVideos := make([]*dto.PlaylistVideo, 0, len(playlist.Videos))
		for _, pv := range playlist.Videos {
			protoVideo := &dto.Video{
				Id:            pv.Video.ID,
				UserId:        pv.Video.UserID,
				Title:         pv.Video.Title,
				Url:           pv.Video.URL,
				ObjectName:    pv.Video.ObjectName,
				EnableComment: pv.Video.EnableComment,
				Visibility:    pv.Video.Visibility,
				CreatedAt:     timestamppb.New(pv.Video.CreatedAt),
				ThumbnailUrl:  pv.Video.ThumbnailUrl,
				PostedAt:      timestamppb.New(pv.Video.PostedAt),
				IsDraft:       pv.Video.IsDraft,
			}
			protoVideos = append(protoVideos, &dto.PlaylistVideo{
				PlaylistId: pv.PlaylistID,
				VideoId:    pv.VideoID,
				Order:      pv.Order,
				Video:      protoVideo,
			})
		}
		protoPlaylists = append(protoPlaylists, &dto.Playlist{
			Id:     playlist.ID,
			UserId: playlist.UserID,
			Title:  playlist.Title,
			Videos: protoVideos,
		})
	}
	return &dto.GetUserPlaylistsResponse{
		Code:      code,
		Message:   message,
		Playlists: protoPlaylists,
	}, nil
}

func (g *GrpcServer) GetPlaylistById(ctx context.Context, in *dto.GetPlaylistByIdRequest) (*dto.GetPlaylistByIdResponse, error) {
	code, message, playlist, err := g.playlistService.GetPlaylistById(ctx, in.PlaylistId)
	if err != nil {
		return &dto.GetPlaylistByIdResponse{
			Code:     code,
			Message:  message,
			Playlist: nil,
		}, err
	}
	var protoPlaylist *dto.Playlist
	if playlist != nil {
		protoVideos := make([]*dto.PlaylistVideo, 0, len(playlist.Videos))
		for _, pv := range playlist.Videos {
			protoVideo := &dto.Video{
				Id:            pv.Video.ID,
				UserId:        pv.Video.UserID,
				Title:         pv.Video.Title,
				Url:           pv.Video.URL,
				ObjectName:    pv.Video.ObjectName,
				EnableComment: pv.Video.EnableComment,
				Visibility:    pv.Video.Visibility,
				CreatedAt:     timestamppb.New(pv.Video.CreatedAt),
				ThumbnailUrl:  pv.Video.ThumbnailUrl,
				PostedAt:      timestamppb.New(pv.Video.PostedAt),
				IsDraft:       pv.Video.IsDraft,
			}
			protoVideos = append(protoVideos, &dto.PlaylistVideo{
				PlaylistId: pv.PlaylistID,
				VideoId:    pv.VideoID,
				Order:      pv.Order,
				Video:      protoVideo,
			})
		}
		protoPlaylist = &dto.Playlist{
			Id:     playlist.ID,
			UserId: playlist.UserID,
			Title:  playlist.Title,
			Videos: protoVideos,
		}
	}
	return &dto.GetPlaylistByIdResponse{
		Code:     code,
		Message:  message,
		Playlist: protoPlaylist,
	}, nil
}

func (g *GrpcServer) UpdatePlaylistTitle(ctx context.Context, in *dto.UpdatePlaylistTitleRequest) (*dto.UpdatePlaylistTitleResponse, error) {
	code, message, err := g.playlistService.UpdatePlaylistTitle(ctx, in.PlaylistId, in.Title)
	return &dto.UpdatePlaylistTitleResponse{
		Code:    code,
		Message: message,
	}, err
}
