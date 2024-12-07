package models

type User struct {
	ID         uint    `gorm:"primaryKey"`
	TelegramID string  `gorm:"uniqueIndex;not null"`
	Username   string  `gorm:"not null"`
	Balance    float64 `gorm:"default:0.0"`
	Ads        []Ad    `gorm:"foreignKey:SellerID"`
	Rating     float64 `gorm:"default:0.0"`
	Language   string  `gorm:"not null;default:'en'"`
	Banned     bool    `gorm:"default:false"`
}
