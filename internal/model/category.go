package model

type Category struct {
	ID   int    `gorm:"autoIncrement;primaryKey" json:"id"`
	Name string `gorm:"not null" json:"name"`
}
