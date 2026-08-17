package domain

import "time"

type WatchHistory struct {
	ID        int32 `gorm:"primaryKey"`
	VideoID   int32 `gorm:"index"`
	UserID    int32 `gorm:"index"`
	Video     Video `gorm:"foreignKey:VideoID"`
	CreatedAt time.Time
}
