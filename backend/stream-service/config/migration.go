package config

import (
	"github.com/Wenev/Survace/stream-service/internal/app/domain"
	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) error {
	return db.AutoMigrate(&domain.StreamChat{})
}
