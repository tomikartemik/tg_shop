package models

type Ad struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Files       string  `gorm:"type:text"`
	SellerID    uint    `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	Approved    bool    `gorm:"default:false"`
}
