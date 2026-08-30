package orderapis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dhirendraj-cmd/cmd/services/orders/models"
)



func GetOrderById(db *sql.DB) http.HandlerFunc{
	fmt.Println("GET ORDER BY ID.......")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idstr := r.URL.Query().Get("id")
		orderID, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, "Invalid ID format. Must be an integer.", http.StatusBadRequest)
			return
		}

		query := `SELECT
						o.status,
						COALESCE(
							json_agg(
								json_build_object(
									'quantity', i.quantity,
									'item_id', i.item_id
								)
							) FILTER (WHERE i.id is NOT NULL), '[]'::json
						) as order_items
				FROM orders o
				LEFT JOIN orderitem i ON o.id = i.order_id
				WHERE o.id=$1
				GROUP BY o.id, o.status;
				`

		var order models.CreateOrderRequest
		var orderItemsRaw []byte // temp container for nested json string

		fmt.Println("Writing query>>>>>>>>>> ")

		err = db.QueryRowContext(r.Context(), query, orderID).Scan(&order.Status, &orderItemsRaw)


		if err!=nil{
			if err == sql.ErrNoRows{
				http.Error(w, "order not found", http.StatusNotFound)
				fmt.Println("No rows for given id is found")
				return 
			}

			http.Error(w, "Error Reading Row for given Id, DB Error", http.StatusInternalServerError)
			fmt.Println("Error reading row", err)
			return 
		}

		fmt.Println("Unmarshalling>>>>>>>>>>>>>> ")

		// unmarshall the raw db
		err = json.Unmarshal(orderItemsRaw, &order.Items)
		if err != nil {
			http.Error(w, "Failed to parse nested data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(order)
		if err!=nil{
			http.Error(w, "Error enoding order", http.StatusInternalServerError)
			fmt.Println("Error while encoding to json.... ", err)
			return
		}

		fmt.Println("END ORDER by ID>>>>>> ")

	}
}
