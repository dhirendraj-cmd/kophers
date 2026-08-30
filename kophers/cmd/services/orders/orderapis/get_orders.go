package orderapis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhirendraj-cmd/cmd/services/orders/models"
)


func GetOrders(db *sql.DB) http.HandlerFunc{
	fmt.Println("Getting orders from db...")

	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM orders")
		if err!=nil{
			http.Error(w, "Error while reading the orders! ", http.StatusBadRequest)
			fmt.Println("Error reading items... ", err)
			return 
		}

		defer rows.Close()

		var allOrders []models.Orders

		for rows.Next(){
			var order models.Orders

			err = rows.Scan(&order.ID, &order.Amount, &order.Status, &order.CreatedAt)
			if err!=nil{
				http.Error(w, "Erorr scanning orders!!!", http.StatusInternalServerError)
				fmt.Println("Error Scaning Items table... ", err)
				return 
			}

			allOrders = append(allOrders, order)
		}
		// check if error occurs while iterating rows
		if err := rows.Err(); err!=nil{
			http.Error(w, "Error Occured while iterating over rows", http.StatusInternalServerError)
			fmt.Println("Error while iterating.... ", err)
			return
		}

		// set content type to application json in header
		w.Header().Set("Content-Type", "application/json")

		// set status code to 200
		w.WriteHeader(http.StatusOK)

		// encode the items to json and write the response body
		err = json.NewEncoder(w).Encode(allOrders)
		if err!=nil{
			http.Error(w, "Error enoding orders", http.StatusInternalServerError)
			fmt.Println("Error while encoding to json.... ", err)
			return
		}

	}
}
