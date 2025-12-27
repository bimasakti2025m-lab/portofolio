package model

type Cart struct {
	Id         int        `json:"id"`
	UserId     string     `json:"user_id"`
	TotalPrice int        `json:"total_price"`
	Items      []CartItem `json:"items"`
}

type CartItem struct {
	Id        int     `json:"id"`
	ProductId int     `json:"product_id"`
	CartId    int     `json:"cart_id"`
	Price     int     `json:"price"`
	Quantity  int     `json:"quantity"`
	Product   Product `json:"product"`
}
