package itemapis

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dhirendraj-cmd/cmd/services/items/models"
)


func CreateItem(db *sql.DB) http.HandlerFunc {
	fmt.Println("POST ITEM API....")

	return func(w http.ResponseWriter, r *http.Request)  {
		
		// create a var to store the item data
		var items models.Items

		// decode data coming from request
		err := json.NewDecoder(r.Body).Decode(&items)

		if err!=nil{
			http.Error(w, "Error decoding item creation request from user", http.StatusBadRequest)
			fmt.Println(err)
			return 
		}


		// execute sql query
		fmt.Println("Executing query to insert data ....... ")

		query := `INSERT INTO items (itemname, description, price, stock, created_at) 
          VALUES ($1, $2, $3, $4, $5) RETURNING id`

		// read the generated ID using queryrow
		err = db.QueryRow(query, items.ItemName, items.Description, items.Price, items.Stock, time.Now()).Scan(&items.ID)

		if err!=nil{
			http.Error(w,"Error while inserting to db ", http.StatusInternalServerError)
			fmt.Println(err)
			return 
		}

		// if insert sucess then just show 201 created msg
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Item Created in DB"))

		fmt.Println("Item Added Successfully..... ")
		
	}
}



