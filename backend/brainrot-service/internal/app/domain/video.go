package domain

import "time"

type Video struct {
	ID            int32 `gorm:"primaryKey"`
	UserID        int32
	ObjectName    string `gorm:"uniqueIndex"`
	URL           string `gorm:"uniqueIndex"`
	Title         string
	ViewCount     int32
	EnableComment bool
	Visibility    string
	ThumbnailUrl  string
	IsDraft       bool
	CreatedAt     time.Time
	PostedAt      time.Time
}
