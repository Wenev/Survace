package main

import (
	"fmt"
	"github.com/Wenev/Survace/brainrot-service/config"
	"github.com/Wenev/Survace/brainrot-service/internal/adapters/inbound/grpc"
	"github.com/Wenev/Survace/brainrot-service/internal/adapters/outbound/cache"
	database "github.com/Wenev/Survace/brainrot-service/internal/adapters/outbound/db"
	"github.com/Wenev/Survace/brainrot-service/internal/adapters/outbound/minio"
	"github.com/Wenev/Survace/brainrot-service/internal/adapters/outbound/repository"
	"github.com/Wenev/Survace/brainrot-service/internal/app"

	brainrotgrpc "github.com/Wenev/Survace/brainrot-service/internal/adapters/inbound/grpc"
	"log"
	"os"
	"strconv"
)

func main() {
	log.Print("Brainrot running PLEASEAANG")
	db, err := database.DatabaseConnect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	err = config.SetupDatabase(db)
	if err != nil {
		panic("failed to run migrations: " + err.Error())
	}
	minioClient, err := minio.MinIOConnect()
	if err != nil {
		log.Fatalf("Failed to connect to minio: %v", err)
		return
	}

	cacheConn := cache.CacheConnection()

	commentRepository := repository.NewCommentRepository(db)
	likeRepository := repository.NewLikeRepository(db)
	videoRepository := repository.NewVideoRepository(db, minioClient)
	watchHistoryRepository := repository.NewWatchHistoryRepository(db)
	replyRepository := repository.NewReplyRepository(db)
	likeCommentRepository := repository.NewLikeCommentRepository(db)
	likeReplyRepository := repository.NewLikeReplyRepository(db)
	playlistRepository := repository.NewPlaylistRepository(db)
	playlistVideoRepository := repository.NewPlaylistVideoRepository(db)
	addressSocial, err := brainrotgrpc.NewSocialClient("social-service:3002")
	if err != nil {
		log.Printf("Failed to connect to galactus: %v", err)
		addressSocial = nil
	}

	if addressSocial == nil {
		log.Printf("Warning: SocialClient is nil, social features will be disabled.")
	}
	videoService := app.NewVideoService(
		videoRepository,
		likeRepository,
		commentRepository,
		watchHistoryRepository,
		playlistRepository,
		playlistVideoRepository,
		addressSocial,
		cacheConn)
	commentService := app.NewCommentService(commentRepository, replyRepository, likeCommentRepository, cacheConn)
	replyService := app.NewReplyService(replyRepository, likeReplyRepository, cacheConn)
	likeService := app.NewLikeService(likeRepository, videoRepository, cacheConn)
	likeCommentService := app.NewLikeCommentService(likeCommentRepository, commentRepository, cacheConn)
	likeReplyService := app.NewLikeReplyService(likeReplyRepository, cacheConn)
	playlistService := app.NewPlaylistService(playlistRepository, playlistVideoRepository, videoRepository, cacheConn)

	portStr := os.Getenv("PORT_BRAINROT")
	port, err := strconv.Atoi(portStr)
	log.Print("PORTTT")
	log.Print(port)
	fmt.Print(port)
	if err != nil || port == 0 {
		log.Printf("Invalid port: %v, defaulting to 3001", err)
		port = 3001
	}
	address, err := brainrotgrpc.NewGalactusClient("galactus:3000")
	if err != nil {
		log.Printf("Failed to connect to galactus: %v", err)
		address = nil
	}
	server := grpc.NewGrpcServer(port, videoService, playlistService, commentService, replyService, likeService, likeCommentService,
		likeReplyService, address)
	defer server.Stop()
	err = server.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
