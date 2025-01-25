package model

type Invoice struct {
	ID         int     `gorm:"autoIncrement;primaryKey" json:"id"`
	TelegramID int     `gorm:"primaryKey;uniqueIndex;not null" json:"telegram_id"`
	Amount     float64 `gorm:"not null" json:"amount"`
	Status     string  `gorm:"not null" json:"status"`
}
