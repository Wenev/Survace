package domain

type LikeComment struct {
	ID        int32 `gorm:"primaryKey"`
	CommentID int32 `gorm:"index"`
	UserID    int32
	Comment   Comment `gorm:"foreignKey:CommentID"`
}
