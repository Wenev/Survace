package config

import (
	"github.com/Acad600-TPA/WEB-WE-251/notification-service/internal/app/domain"
	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.Notification{})
	if err != nil {
		return err
	}
	return nil
}
