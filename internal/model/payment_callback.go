package model

type PaymentCallback struct {
	Status  string `json:"status"`
	OrderID int    `json:"order_id"`
}
