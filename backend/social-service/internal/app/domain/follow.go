package domain

import "time"

type Follow struct {
	ID         int32 `gorm:"primaryKey"`
	FollowerID int32 `gorm:"index"`
	FolloweeID int32 `gorm:"index"`
	CreatedAt  time.Time
}
