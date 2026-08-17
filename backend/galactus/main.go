package main

import (
	"github.com/Acad600-TPA/WEB-WE-251/galactus/config"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/adapters/inbound/grpc"
	cache2 "github.com/Acad600-TPA/WEB-WE-251/galactus/internal/adapters/outbound/cache"
	database "github.com/Acad600-TPA/WEB-WE-251/galactus/internal/adapters/outbound/db"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/adapters/outbound/minio"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/adapters/outbound/repository"
	"github.com/Acad600-TPA/WEB-WE-251/galactus/internal/app"
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
	minioClient, err := minio.MinIOConnect()
	if err != nil {
		log.Fatalf("Failed to connect to minio: %v", err)
		return
	}
	cache := cache2.CacheConnection()
	userRepository := repository.NewUserRepository(db, cache, minioClient)
	userService := app.NewUserService(userRepository)
	settingRepository := repository.NewSettingRepository(db)
	settingService := app.NewSettingService(settingRepository, userRepository, cache)
	port, err := strconv.Atoi(os.Getenv("PORT_GALACTUS"))
	if err != nil {
		log.Printf("Invalid port: %v", err)
		port = 3000
	}

	server := grpc.NewGrpcServer(port, userService, settingService)
	defer server.Stop()
	server.Start()
	err = server.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
