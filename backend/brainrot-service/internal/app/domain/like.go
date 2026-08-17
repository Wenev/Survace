package domain

type Like struct {
	ID      int32 `gorm:"primaryKey"`
	VideoID int32 `gorm:"index"`
	UserID  int32
	Video   Video `gorm:"foreignKey:VideoID"`
}
