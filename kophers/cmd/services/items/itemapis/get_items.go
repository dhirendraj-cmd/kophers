package itemapis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dhirendraj-cmd/cmd/services/items/models"
)


func GetItems(db *sql.DB) http.HandlerFunc {
	fmt.Println("Getting Items from db")

	return func(w http.ResponseWriter, r *http.Request) {

		// get all items
		rows, err := db.Query("SELECT * FROM items")

		if err!=nil{
			http.Error(w, "Error while reading the items! ", http.StatusBadRequest)
			fmt.Println("Error reading items... ", err)
			return 
		}

		defer rows.Close()

		var allItems []models.Items

		// iterate over rows and store each item
		for rows.Next(){

			var item models.Items

			// scan rows
			err := rows.Scan(&item.ID, &item.ItemName, &item.Description, &item.Price, &item.Stock, &item.CreatedAt)
			if err!=nil{
				http.Error(w, "Erorr scanning items!!!", http.StatusInternalServerError)
				fmt.Println("Error Scaning Items table... ", err)
				return 
			}

			allItems = append(allItems, item)
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
		err = json.NewEncoder(w).Encode(allItems)
		if err!=nil{
			http.Error(w, "Error enoding items", http.StatusInternalServerError)
			fmt.Println("Error while encoding to json.... ", err)
			return
		}
	}

}
