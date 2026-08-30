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
	fmt.Println("Get order by id.... .")

	return func(w http.ResponseWriter, r *http.Request) {
		idstr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			http.Error(w, "Invalid ID format. Must be an integer.", http.StatusBadRequest)
			return
		}

		var order models.Orders
		query := `SELECT id, amount, status, created_at FROM items WHERE id=$1`

		err = db.QueryRowContext(r.Context(), query, id).Scan(&order.ID, &order.Amount, &order.Status, &order.CreatedAt)

		if err!=nil{
			if err == sql.ErrNoRows{
				http.Error(w, "order not found", http.StatusNotFound)
				fmt.Println("No rows for given id is found")
				return 
			}

			http.Error(w, "Error Reading Row for given Id", http.StatusInternalServerError)
			fmt.Println("Error reading row")
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


	}
}
