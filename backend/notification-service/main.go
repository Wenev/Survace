package main

import (
	"github.com/Wenev/Survace/notification-service/config"
	"github.com/Wenev/Survace/notification-service/internal/adapters/inbound/grpc"
	"github.com/Wenev/Survace/notification-service/internal/adapters/outbound/cache"
	database "github.com/Wenev/Survace/notification-service/internal/adapters/outbound/db"
	"github.com/Wenev/Survace/notification-service/internal/adapters/outbound/repository"
	"github.com/Wenev/Survace/notification-service/internal/app"
	"log"
	"os"
	"strconv"
)

func main() {
	db, err := database.DatabaseConnect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	if err := config.SetupDatabase(db); err != nil {
		panic("failed to run migrations: " + err.Error())
	}
	if err != nil {
		log.Fatalf("Failed to connect to minio: %v", err)
		return
	}
	cache := cache.CacheConnection()
	userRepository := repository.NewNotificationRepository(db)
	userService := app.NewNotificationService(userRepository, cache)
	port, err := strconv.Atoi(os.Getenv("PORT_NOTIF"))
	if err != nil {
		log.Printf("Invalid port: %v", err)
		port = 3000
	}

	server := grpc.NewGrpcServer(port, userService)
	defer server.Stop()
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
