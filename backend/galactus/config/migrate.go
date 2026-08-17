package config

import (
	"github.com/Wenev/Survace/galactus/internal/app/domain"
	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		return err
	}
	return db.AutoMigrate(&domain.Setting{})
}
