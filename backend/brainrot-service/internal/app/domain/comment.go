package domain

type Comment struct {
	ID      int32 `gorm:"primaryKey"`
	Content string
	VideoID int32 `gorm:"index"`
	UserID  int32
	Video   Video `gorm:"foreignKey:VideoID"`
}
