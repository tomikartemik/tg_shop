package model

type Category struct {
	ID   int    `gorm:"autoIncrement;primaryKey" `
	Name string `gorm:"not null" json:"name"`
}
