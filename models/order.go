package models

import "time"

type TransactionItem struct {
	Title string  `json:"title"`
	Qty   int     `json:"qty"`
	Size  string  `json:"size,omitempty"` 
}

type Transaction struct {
	ID            int64             `json:"id"`
	NoOrders      string            `json:"no_orders"`
	CreatedAt     time.Time         `json:"created_at"`
	StatusName    string            `json:"status_name"`
	Total         float64           `json:"total"`
	UserFullname  string            `json:"user_fullname"`
	UserAddress   string            `json:"user_address"`
	UserPhone     string            `json:"user_phone"`
	PaymentMethod string            `json:"payment_method"`
	ShippingName  string            `json:"shipping_name"`
	OrderItems    []TransactionItem `json:"order_items"`
}
