package model

type Ad struct {
	ID          int     `gorm:"primaryKey" json:"id"`
	Title       string  `gorm:"not null" json:"title"`
	Description string  `gorm:"not null" json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Files       string  `gorm:"type:text" json:"files_url"`
	PhotoURL    string  `gorm:"type:text" json:"photo_url"`
	CategoryID  int     `gorm:"not null" json:"category_id"`
	SellerID    int     `gorm:"not null" json:"seller_id"`
	Stock       int     `gorm:"not null" json:"stock"`
	Status      string  `gorm:"default:Rejected"`
}

type AdInfo struct {
	ID                 int     `json:"id"`
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	Price              float64 `json:"price"`
	Files              string  `json:"files_url"`
	PhotoURL           string  `json:"photo_url"`
	CategoryID         int     `json:"category_id"`
	CategoryName       string  `json:"category_name"`
	SellerID           int     `json:"seller_id"`
	SellerName         string  `json:"seller_name"`
	SellerRating       float64 `json:"seller_rating"`
	SellerReviewNumber int     `json:"seller_review_number"`
	Stock              int     `json:"stock"`
}

type AdShortInfo struct {
	ID                 int     `json:"id"`
	Title              string  `json:"title"`
	Price              float64 `json:"price"`
	PhotoURL           string  `json:"photo_url"`
	CategoryID         int     `json:"category_id"`
	CategoryName       string  `json:"category_name"`
	SellerID           int     `json:"seller_id"`
	SellerName         string  `json:"seller_name"`
	SellerRating       float64 `json:"seller_rating"`
	SellerReviewNumber int     `json:"seller_review_number"`
	Stock              int     `json:"stock"`
}
