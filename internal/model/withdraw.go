package model

type PayoutRequest struct {
	ID         int     `gorm:"primaryKey"`
	TelegramID int     `gorm:"index"`
	Amount     float64 `gorm:"not null"`
	Status     string  `gorm:"default:'pending'"`
}
