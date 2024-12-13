package model

import "time"

type User struct {
	TelegramID    int       `gorm:"primaryKey;uniqueIndex;not null"`
	Username      string    `gorm:"not null;default:'default_user'"`
	IsAdmin       bool      `gorm:"not null;default:false"`
	Balance       float64   `gorm:"default:0.0"`
	Ads           []Ad      `gorm:"foreignKey:SellerID"`
	Purchased     []Ad      `gorm:"many2many:user_purchased_ads"`
	Rating        float64   `gorm:"default:0.0"`
	ReviewNumber  int       `gorm:"default:0"`
	Banned        bool      `gorm:"default:false"`
	IsPremium     bool      `gorm:"not null;default:false"`
	ExpirePremium time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type UserAsSeller struct {
	TelegramID   int           `json:"telegram_id"`
	Username     string        `json:"username"`
	Ads          []AdShortInfo `json:"ads"`
	Rating       float64       `json:"rating"`
	ReviewNumber int           `json:"reviewNumber"`
}
