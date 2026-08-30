package orderapis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"

	"github.com/dhirendraj-cmd/cmd/services/orders/models"
)



func CreateOrder(db *sql.DB) http.HandlerFunc{
	fmt.Println("POST Order Creation.... ")

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Starting Creation process...>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		var req models.CreateOrderRequest
		var orderItems []models.OrderItem
		var orderId, amounts int

		fmt.Println("decoding process...>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		err := json.NewDecoder(r.Body).Decode(&req)
		if err!=nil{
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			fmt.Println(err)
			return 
		}

		fmt.Println("transaction process...>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		// begin database transaction
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		fmt.Println("Rollback process...>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		// rollback
		defer tx.Rollback()

		fmt.Println("fetch prices and stock process...>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		// fetch prices
		for _, item := range req.Items{
			var prices, stock int
			item_query := `SELECT price, stock FROM items where id=$1`
			err = tx.QueryRowContext(r.Context(), item_query, item.ItemId).Scan(&prices, &stock)
			if err != nil {
				http.Error(w, "Order not found", http.StatusBadRequest)
				fmt.Println("Prices Query error>>>>>>>>>>>>>>>>>>> ", err)
				return
			}

			fmt.Println("ITem id stock is ", stock, item.ItemId)

			if (stock >= item.Quantity){
				fmt.Println("In STOCK>>>>>>>>>>>>>>")
				amounts += prices*item.Quantity
				
				orderItems = append(orderItems, models.OrderItem{
					ItemId:   item.ItemId,
					Quantity: item.Quantity,
					Price:    prices,
				})
				
				stock = stock-item.Quantity
				stock_query := `UPDATE items SET stock=$1 WHERE id=$2`
				res, err := tx.ExecContext(r.Context(), stock_query, stock, item.ItemId)
				if err!=nil{
					http.Error(w, "Unable to update Stock", http.StatusBadRequest)
					fmt.Println("Unable to update stock due to>>> ", err)
					return 
				}
				affectedRows, err := res.RowsAffected()
				if err!=nil{
					http.Error(w, "Error verifying update", http.StatusInternalServerError)
					fmt.Println("Error verifying update due to>>> ", err)
					fmt.Println()
					return
				}
				if affectedRows==0{
					http.Error(w, "Record not found", http.StatusNotFound)
					fmt.Println("Record not found>>>>>>>>")
					fmt.Println()
					return
				}

				} else {
					w.Write([]byte("Order Cannot be created, because item in stock is less than ordered quantity"))
					// fmt.Sprintf("Item in Stock is less than Quantity, only %v items are present in stock", stock)
					fmt.Println("Not In STOCK>>>>>>>>>>>>>>")
				return 
			} 

		}

		// create order insert parent record
		order_query := `INSERT INTO orders (amount, status) VALUES ($1, $2) RETURNING id`
		err = tx.QueryRowContext(r.Context(), order_query, amounts, req.Status).Scan(&orderId)
		if err!=nil{
			http.Error(w,"Error while inserting to db ", http.StatusInternalServerError)
			fmt.Println("Insertion Query Error for orders>>>>>>>>>>> ", err)
			return 
		}

		// insert child record
		order_item_query := `INSERT INTO orderitem (price, quantity, order_id, item_id) VALUES ($1, $2, $3, $4)`
		for _, item:=range orderItems{
			_, err = tx.ExecContext(r.Context(), order_item_query, item.Price, item.Quantity, orderId, item.ItemId)

			if err != nil {
				http.Error(w, "Error while inserting order items", http.StatusInternalServerError)
				fmt.Println("Insertion Query Error for order Items>>>>>>>>>>> ", err)
				return
			}

		}

		err = tx.Commit()
		if err != nil {
            http.Error(w,"Error committing transaction", http.StatusInternalServerError)
			fmt.Println("Error while commiting for orders>>>>>>>>>>> ", err)
            return
        }

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Order Created"))
		
		fmt.Println("Order Created Successfully..... ")

	}
}

