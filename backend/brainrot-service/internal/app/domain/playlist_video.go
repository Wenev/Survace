package domain

type PlaylistVideo struct {
	ID         int32 `gorm:"primaryKey"`
	PlaylistID int32 `gorm:"index"`
	VideoID    int32 `gorm:"index"`
	Order      int32
	Video      Video    `gorm:"foreignKey:VideoID"`
	Playlist   Playlist `gorm:"foreignKey:PlaylistID"`
}
