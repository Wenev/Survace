package main

import (
	"log"
	"os"
	"strconv"

	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/config"
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/adapters/inbound/grpc"
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/adapters/outbound/cache"
	database "github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/adapters/outbound/db"
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/adapters/outbound/repository"
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/app"
)

func main() {
	log.Print("running sosial")
	db, err := database.DatabaseConnect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	if err := config.SetupDatabase(db); err != nil {
		panic("failed to run migrations: " + err.Error())
	}
	cacheConn := cache.CacheConnection()
	chatRepo := repository.NewChatRepository(db)
	followRepo := repository.NewFollowRepository(db)
	chatService := app.NewChatService(chatRepo, cacheConn)
	followService := app.NewFollowService(followRepo)
	port, err := strconv.Atoi(os.Getenv("PORT_SOCIAL"))
	if err != nil {
		log.Printf("Invalid port: %v", err)
		port = 3000
	}
	address, err := grpc.NewGalactusClient("galactus:3000", "notification-service:3004")
	if err != nil {
		log.Printf("Invalid port: %v", err)
		port = 3000
	}
	server := grpc.NewGrpcServer(port, chatService, followService, address)
	defer server.Stop()
	err = server.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
