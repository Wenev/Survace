package domain

import "time"

type Notification struct {
	ID        int32  `gorm:"primaryKey"`
	UserID    int32  `gorm:"index"`
	Type      string `gorm:"not null"`
	Content   string
	Read      bool `gorm:"default:false"`
	CreatedAt time.Time
}
