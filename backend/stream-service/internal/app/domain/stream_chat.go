package domain

import "time"

type StreamChat struct {
	ID        int32 `gorm:"primaryKey"`
	SenderID  int32 `gorm:"index"`
	CallID    int32 `gorm:"index"`
	Content   string
	CreatedAt time.Time
}
