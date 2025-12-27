package model

type OrderItem struct {
	ID            int     `json:"id"`
	OrderID       int     `json:"order_id"`
	ProductID     int     `json:"product_id"`
	Quantity      int     `json:"quantity"`
	PriceSnapshot float64 `json:"price_snapshot"`
}
