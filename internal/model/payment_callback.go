package model

type PaymentCallback struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}
