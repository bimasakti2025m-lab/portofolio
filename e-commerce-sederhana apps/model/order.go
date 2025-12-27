package model

type Order struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Total float64 `json:"total"`
	StatusPesanan string `json:"status_pesanan"`
	TransactionIDMidtrans string `json:"transaction_id_midtrans"`
}

