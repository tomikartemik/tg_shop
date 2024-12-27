package model

import "time"

type User struct {
	TelegramID     int     `gorm:"primaryKey;uniqueIndex;not null" json:"telegram_id"`
	Username       string  `gorm:"not null;unique" json:"username"`
	PhotoURL       string  `gorm:"not null;default:'uploads/default_hell_ava.jpg" json:"photo_url"`
	IsAdmin        bool    `gorm:"not null;default:false" json:"is_admin"`
	Balance        float64 `gorm:"default:0.0" json:"balance"`
	Ads            []Ad    `gorm:"foreignKey:SellerID" json:"ads"`
	Purchased      []Ad    `gorm:"many2many:user_purchased_ads" json:"purchased"`
	PurchaseAmount float64 `gorm:"default:0.0" json:"purchase_amount"`
	//PurchaseAmount float64 `json:"purchase_amount"`
	Rating       float64 `gorm:"default:0.0" json:"rating"`
	ReviewNumber int     `gorm:"default:0" json:"review_number"`
	SalesAmount  float64 `gorm:"default:0.0" json:"sales_amount"`
	//SalesAmount   float64   `json:"sales_amount"`
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

type UserInfo struct {
	TelegramID     int           `json:"telegram_id"`
	Username       string        `json:"username"`
	PhotoURL       string        `json:"photo_url"`
	Balance        float64       `json:"balance"`
	Ads            []AdShortInfo `json:"ads"`
	Purchased      []AdShortInfo `json:"purchased"`
	PurchaseAmount float64       `json:"purchase_amount"`
	Rating         float64       `json:"rating"`
	ReviewNumber   int           `json:"review_number"`
	SalesAmount    float64       `json:"sales_amount"`
}
