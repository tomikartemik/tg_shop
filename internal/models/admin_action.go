package models

type AdminAction struct {
	ID        uint   `gorm:"primaryKey"`
	Action    string `gorm:"not null"`
	Timestamp string `gorm:"autoCreateTime"`
	UserID    uint   `gorm:"not null"`
}
