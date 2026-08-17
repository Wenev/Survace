package domain

type Playlist struct {
	ID     int32 `gorm:"primaryKey"`
	Title  string
	UserID int32
	Videos []PlaylistVideo `gorm:"foreignKey:PlaylistID"`
}
