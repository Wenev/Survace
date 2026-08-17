package domain

import "time"

type Notification struct {
	ID        int32  `gorm:"primaryKey"`
	UserID    int32  `gorm:"index"`    // The user who receives the notification
	Type      string `gorm:"not null"` // Validate in code, not as enum in DB
	Content   string
	Read      bool `gorm:"default:false"`
	CreatedAt time.Time
}
