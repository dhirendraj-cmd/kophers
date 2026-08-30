package itemapis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dhirendraj-cmd/cmd/services/items/models"
)


func GetItemById(db *sql.DB) http.HandlerFunc{
	fmt.Println("Clling get item by id... ")

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID format. Must be an integer.", http.StatusBadRequest)
			return
		}

		var item models.Items

		query := `SELECT id, itemname, description, price, stock, created_at FROM items WHERE id=$1`

		err = db.QueryRowContext(r.Context(), query, id).Scan(&item.ID, &item.ItemName, &item.Description, &item.Price, &item.Stock, &item.CreatedAt)

		if err!=nil{
			if err == sql.ErrNoRows{
				http.Error(w, "item not found", http.StatusNotFound)
				fmt.Println("No rows for given id is found")
				return 
			}

			http.Error(w, "Error Reading Row for given Id", http.StatusInternalServerError)
			fmt.Println("Error reading row")
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(item)
		if err!=nil{
			http.Error(w, "Error enoding item", http.StatusInternalServerError)
			fmt.Println("Error while encoding to json.... ", err)
			return
		}


	}
}
