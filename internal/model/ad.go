package model

type Ad struct {
	ID          int     `gorm:"primaryKey"`
	Title       string  `gorm:"not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Files       string  `gorm:"type:text"`
	PhotoURL    string  `gorm:"type:text"`
	CategoryID  int     `gorm:"not null"`
	SellerID    int     `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	Approved    bool    `gorm:"default:false"`
}

type AdInfo struct {
	ID                 int     `json:"id"`
	Title              string  `json:"title"`
	Description        string  `json:"description"`
	Price              float64 `json:"price"`
	Files              string  `json:"files"`
	PhotoURL           string  `gorm:"type:text"`
	CategoryID         int     `json:"categoryID"`
	CategoryName       string  `json:"categoryName"`
	SellerID           int     `json:"seller_id"`
	SellerName         string  `json:"seller_name"`
	SellerRating       float64 `json:"seller_rating"`
	SellerReviewNumber int     `json:"seller_review_number"`
	Stock              int     `gorm:"not null"`
}

type AdShortInfo struct {
	ID                 int     `json:"id"`
	Title              string  `json:"title"`
	Price              float64 `json:"price"`
	PhotoURL           string  `gorm:"type:text"`
	CategoryID         int     `json:"categoryID"`
	CategoryName       string  `json:"categoryName"`
	SellerID           int     `json:"seller_id"`
	SellerName         string  `json:"seller_name"`
	SellerRating       float64 `json:"seller_rating"`
	SellerReviewNumber int     `json:"seller_review_number"`
	Stock              int     `gorm:"not null"`
}
