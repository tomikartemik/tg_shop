package model

type Ad struct {
	ID          uint    `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Files       string  `gorm:"type:text"`
	CategoryID  int     `gorm:"not null"`
	SellerID    string  `gorm:"uniqueIndex;not null"`
	Stock       int     `gorm:"not null"`
	Approved    bool    `gorm:"default:false"`
}
