package models

import "time"


type Items struct{
	ID			int 		`json:"id"`
	ItemName	string 		`json:"itemname"`
	Description	string 		`json:"description"`
	Price		int 		`json:"price"`
	Stock		int 		`json:"stock"`
	CreatedAt	time.Time 	`json:"created_at"`
}

