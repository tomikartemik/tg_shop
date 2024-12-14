package model

import "time"

type User struct {
	TelegramID    int       `gorm:"primaryKey;uniqueIndex;not null" json:"telegram_id"`
	Username      string    `gorm:"not null;default:'default_user'" json:"username"`
	PhotoURL      string    `gorm:"not null;default:'uploads/default_hell_ava.jpg" json:"photo_url"`
	IsAdmin       bool      `gorm:"not null;default:false" json:"is_admin"`
	Balance       float64   `gorm:"default:0.0" json:"balance"`
	Ads           []Ad      `gorm:"foreignKey:SellerID" json:"ads"`
	Purchased     []Ad      `gorm:"many2many:user_purchased_ads" json:"purchased"`
	Rating        float64   `gorm:"default:0.0" json:"rating"`
	ReviewNumber  int       `gorm:"default:0" json:"review_number"`
	Banned        bool      `gorm:"default:false" json:"banned"`
	IsPremium     bool      `gorm:"not null;default:false" json:"is_premium"`
	ExpirePremium time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"expire_premium"`
}

type UserAsSeller struct {
	TelegramID   int           `json:"telegram_id"`
	Username     string        `json:"username"`
	PhotoURL     string        `json:"photo_url"`
	Ads          []AdShortInfo `json:"ads"`
	Rating       float64       `json:"rating"`
	ReviewNumber int           `json:"review_number"`
}
