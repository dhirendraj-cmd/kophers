package orderapis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhirendraj-cmd/cmd/services/orders/models"
)



func CreateOrder(db *sql.DB) http.HandlerFunc{
	fmt.Println("Order Creation.... ")

	return func(w http.ResponseWriter, r *http.Request) {

		var req models.CreateOrderRequest
		var orderItems []models.OrderItem
		var orderId, amounts int

		err := json.NewDecoder(r.Body).Decode(&req)
		if err!=nil{
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			fmt.Println(err)
			return 
		}

		// begin database transaction
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// rollback
		defer tx.Rollback()

		// fetch prices
		for _, item := range req.Items{
			var prices int
			item_prices := `SELECT price FROM items where id=$1`
			err = tx.QueryRowContext(r.Context(), item_prices, item.ItemId).Scan(&prices)
			if err != nil {
				http.Error(w, "Item not found", http.StatusBadRequest)
				return
			}

			amounts += prices*item.Quantity

			orderItems = append(orderItems, models.OrderItem{
				ItemId:   item.ItemId,
				Quantity: item.Quantity,
				Price:    prices,
			})
		}

		// create order insert parent record
		order_query := `INSERT INTO orders (amount, status) VALUES ($1, $2) RETURNING id`
		err = tx.QueryRowContext(r.Context(), order_query, amounts, req.Status).Scan(&orderId)
		if err!=nil{
			http.Error(w,"Error while inserting to db ", http.StatusInternalServerError)
			fmt.Println(err)
			return 
		}

		// insert child record
		order_item_query := `INSERT INTO orderitem (price, quantity, order_id, item_id) VALUES ($1, $2, $3, $4)`
		for _, item:=range orderItems{
			_, err = tx.ExecContext(r.Context(), order_item_query, item.Price, item.Quantity, orderId, item.ItemId)

			if err != nil {
				http.Error(w, "Error while inserting order items", http.StatusInternalServerError)
				return
			}

		}

		err = tx.Commit()
		if err != nil {
            http.Error(w,"Error committing transaction", http.StatusInternalServerError)
            return
        }

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Item Created in DB"))

		fmt.Println("Order Created Successfully..... ")

	}
}

