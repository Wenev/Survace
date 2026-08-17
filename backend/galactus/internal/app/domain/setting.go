package domain

//go:generate stringer -type=ChatRestrictionType

type ChatRestrictionType string

const (
	ChatRestrictionFriends  ChatRestrictionType = "friends"
	ChatRestrictionEveryone ChatRestrictionType = "everyone"
	ChatRestrictionNone     ChatRestrictionType = "none"
)

type Setting struct {
	ID                      int32 `gorm:"primaryKey"`
	UserID                  int32 `gorm:"uniqueIndex"`
	ChatRestriction         ChatRestrictionType
	Private                 bool `gorm:"default:false"`
	NewFollowerNotification bool `gorm:"default:true"`
	MessageNotification     bool `gorm:"default:true"`
	MentionsNotification    bool `gorm:"default:true"`
	LikeTabVisibility       bool `gorm:"default:true"`
}
