package model

type PurchaseRequest struct {
	UserID int `json:"telegram_id"`
	AdID   int `json:"ad_id"`
}
