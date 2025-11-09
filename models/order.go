package models

import "time"

type Order struct {
	ID               int64           `json:"id"`
	NoOrders         string          `json:"no_orders"`
	Total            float64         `json:"total"`
	StatusID         int64           `json:"status_id"`
	StatusName       string          `json:"status_name"`
	UserID           int64           `json:"user_id"`
	UserName         string          `json:"user_name"`
	PaymentMethodID  int64           `json:"payment_method_id"`
	PaymentMethod    string          `json:"payment_method"`
	ShippingName     string          `json:"shipping_name"`
	CreatedAt        time.Time       `json:"created_at"`
	OrderItems       []OrderItemResp `json:"order_items,omitempty"`
}

type OrderItemResp struct {
	Title string `json:"title"`
	Size  string `json:"size"`
	Qty   int    `json:"qty"`
}
