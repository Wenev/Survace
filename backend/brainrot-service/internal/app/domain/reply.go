package domain

type Reply struct {
	ID        int32 `gorm:"primaryKey"`
	Content   string
	CommentID int32 `gorm:"index"`
	UserID    int32
	Comment   Comment `gorm:"foreignKey:CommentID"`
}
