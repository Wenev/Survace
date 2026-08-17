package domain

import "time"

type User struct {
	ID              int32  `gorm:"primaryKey"`
	Username        string `gorm:"uniqueIndex"`
	Email           string `gorm:"uniqueIndex"`
	Password        string
	IsEmailVerified bool `gorm:"default:false"`
	DateOfBirth     time.Time
	AvatarURL       string
}
