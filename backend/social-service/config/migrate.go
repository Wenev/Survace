package config

import (
	"github.com/Acad600-TPA/WEB-WE-251/backend/social-service/internal/app/domain"
	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.Follow{})
	if err != nil {
		return err
	}
	return db.AutoMigrate(&domain.Message{})
}
