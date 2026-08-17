package config

import (
	"github.com/Acad600-TPA/WEB-WE-251/brainrot-service/internal/app/domain"
	"gorm.io/gorm"
)

func SetupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(&domain.Video{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.Playlist{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.PlaylistVideo{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.Comment{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.Like{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.WatchHistory{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.LikeComment{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.LikeReply{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&domain.Reply{})
	if err != nil {
		return err
	}
	return nil
}
