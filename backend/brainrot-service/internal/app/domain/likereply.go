package domain

type LikeReply struct {
	ID      int32 `gorm:"primaryKey"`
	ReplyID int32 `gorm:"index"`
	UserID  int32
	Reply   Reply `gorm:"foreignKey:ReplyID"`
}
