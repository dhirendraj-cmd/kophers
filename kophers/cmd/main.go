package main

import (
	"fmt"
	"net/http"

	"github.com/dhirendraj-cmd/cmd/databases"
	"github.com/dhirendraj-cmd/cmd/services/items/itemapis"
	"github.com/dhirendraj-cmd/cmd/services/orders/orderapis"
)



func main(){
	fmt.Println("KOPHERS")
	
	db := databases.Connection()
	defer db.Close()


	// hanling the items apis
	http.HandleFunc("/api/itm/create_item", itemapis.CreateItem(db))
	http.HandleFunc("/api/itm/items", itemapis.GetItems(db))
	http.HandleFunc("/api/itm/item", itemapis.GetItemById(db))

	// order apis
	http.HandleFunc("/api/ord/create", orderapis.CreateOrder(db))
	http.HandleFunc("/api/ord/orders", orderapis.GetOrders(db))
	http.HandleFunc("/api/ord/order", orderapis.GetOrderById(db))

	err := http.ListenAndServe(":3000", nil)
	if err!=nil{
		fmt.Println("Error while starting HTTP Server: ", err)
		return
	}

}
