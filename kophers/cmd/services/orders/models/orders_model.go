package models

import "time"


type Orders struct{
	ID			int 		`json:"id"`
	Amount		int 		`json:"amount"`
	Status 		string 		`json:"status"`
	CreatedAt	time.Time 	`json:"created_at"`
}

type OrderItem struct{
	ID			int `json:"id"`
	Price		int `json:"price"`
	Quantity	int `json:"quantity"`
	OrderId		int `json:"order_id"`
	ItemId		int `json:"item_id"`
}

// child
type OrderItemRequest struct {
	Quantity int `json:"quantity"`
	ItemId   int `json:"item_id"`
}

// parent
type CreateOrderRequest struct {
    Status string `json:"status"`
    Items  []OrderItemRequest `json:"items"`
}
