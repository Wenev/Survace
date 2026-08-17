package config

import (
	"github.com/Wenev/Survace/notification-service/internal/app/domain"
	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.Notification{})
	if err != nil {
		return err
	}
	return nil
}
