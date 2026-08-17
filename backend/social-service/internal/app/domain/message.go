package domain

import "time"

type Message struct {
	ID         int32 `gorm:"primaryKey"`
	SenderID   int32 `gorm:"index"`
	ReceiverID int32 `gorm:"index"`
	Content    string
	CreatedAt  time.Time
}
