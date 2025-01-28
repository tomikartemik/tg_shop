package model

import "time"

type Earning struct {
	ID          int        `gorm:"autoIncrement;primaryKey" json:"id"`
	BuyerID     int        `gorm:"not null" json:"buyerId"`
	SellerID    int        `gorm:"not null" json:"sellerId"`
	Amount      float64    `gorm:"not null" json:"amount"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Status      string     `gorm:"not null" json:"status"`
	ProcessedAt *time.Time `gorm:"default:null" json:"processed_at"`
}
