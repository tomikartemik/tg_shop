package model

type User struct {
	TelegramID   int     `gorm:"primaryKey;uniqueIndex;not null"`
	Username     string  `gorm:"not null;default:'default_user'"`
	IsAdmin      bool    `gorm:"not null;default:false"`
	Balance      float64 `gorm:"default:0.0"`
	Ads          []Ad    `gorm:"foreignKey:SellerID"`
	Purchased    []Ad    `gorm:"many2many:user_purchased_ads"`
	Rating       float64 `gorm:"default:0.0"`
	ReviewNumber int     `gorm:"default:0"`
	Language     string  `gorm:"not null;default:en"`
	Banned       bool    `gorm:"default:false"`
}
