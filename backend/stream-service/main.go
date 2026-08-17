package main

import (
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/adapters/outbound/cache"
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/adapters/outbound/db"
	"log"
	"os"
	"strconv"

	"github.com/Acad600-TPA/WEB-WE-251/stream-service/config"
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/adapters/inbound/grpc"
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/adapters/outbound/repository"
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	log.Print("running stream-service")
	db, err := db.DatabaseConnect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	if err := config.SetupDatabase(db); err != nil {
		panic("failed to run migrations: " + err.Error())
	}
	cacheConn := cache.CacheConnection()
	chatRepo := repository.NewStreamChatRepository(db)
	chatService := app.NewStreamChatService(chatRepo)

	// Create stream service with API credentials
	streamService := app.NewStreamService(cacheConn)

	port, err := strconv.Atoi(os.Getenv("PORT_STREAM"))
	if err != nil {
		log.Printf("Invalid port: %v", err)
		port = 3000
	}

	server := grpc.NewGrpcServer(port, streamService, chatService)
	defer server.Stop()
	err = server.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
