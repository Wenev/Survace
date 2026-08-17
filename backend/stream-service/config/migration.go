package config

import (
	"github.com/Acad600-TPA/WEB-WE-251/stream-service/internal/app/domain"
	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) error {
	return db.AutoMigrate(&domain.StreamChat{})
}
