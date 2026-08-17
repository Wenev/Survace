package db

import (
	"fmt"
	"github.com/Wenev/Survace/backend/social-service/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DatabaseConnect() (*gorm.DB, error) {
	dbURL := config.DbURLFromEnv()
	fmt.Print(dbURL)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	return db, nil
}
