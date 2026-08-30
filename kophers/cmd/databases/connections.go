package databases

import (
	"database/sql"
	"fmt"
	_"github.com/lib/pq"
)


const (
	host 		= "<host>"
	port 		= "<port-no>"
	user 		= "<username>"
	password 	= "<youpwd>"
	dbname 	 	= "<dbname>"
)


func CreateTableSchema(db *sql.DB) error{
	// create table if not exists
	item_table := `CREATE TABLE IF NOT EXISTS items(
						id	SERIAL PRIMARY KEY,			
						itemname VARCHAR(255) NOT NULL,
						description TEXT,
						price	INT NOT NULL CHECK (price > 0),
						stock	INT NOT NULL CHECK (stock >= 0),
						created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
					);`

	if _, err := db.Exec(item_table); err !=nil{
		return err
	}

	order_table := `CREATE TABLE IF NOT EXISTS orders(
						id	SERIAL PRIMARY KEY,
						amount INT NOT NULL DEFAULT 0,
						status VARCHAR(255) NOT NULL DEFAULT 'CREATED',
						created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
					);`

	if _, err := db.Exec(order_table); err !=nil{
		return err
	}


	order_item_table := `CREATE TABLE IF NOT EXISTS orderitem(
							id	SERIAL PRIMARY KEY,

							price	INT NOT NULL CHECK (price > 0),
							quantity INT NOT NULL CHECK (quantity > 0),

							order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
							item_id INT NOT NULL REFERENCES items(id) ON DELETE CASCADE
					);`

	if _, err := db.Exec(order_item_table); err !=nil{
		return err
	}
		
	return nil

}


func Connection() *sql.DB{
	fmt.Println("Trying to connect DB....... ")

	// establish connection string to connect with db
	fmt.Println("creating connection string..... ")
	connStr := fmt.Sprintf("host=%s\n port=%d\n user=%s\n password=%s\n dbname=%s\n sslmode=disable", host, port, user, password, dbname)

	// opening db to connect
	fmt.Println("Opening DB to connect using sql.....")
	db, err := sql.Open("postgres", connStr)

	if err != nil{
		panic(err)
	}

	// defer db.Close()

	// run table checking
	err = CreateTableSchema(db)
	if err!=nil{
		fmt.Println("Error while creating items table....")
	}
	fmt.Println("Table created")

	// setting db pool connections
	fmt.Println("Setting Pooling connection..... ")

	// set max open connections
	db.SetMaxOpenConns(25)

	// setting max idle conn
	db.SetMaxIdleConns(25)

	// setting conn life time
	db.SetConnMaxLifetime(0)


	// testing connection
	fmt.Println("pinging and testing conn ......")
	err = db.Ping()
	if err!=nil{
		panic(err)
	}

	fmt.Println("Successfully conneted.......")

	return db

}


