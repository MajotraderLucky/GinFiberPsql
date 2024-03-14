package checkdb

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectAndCheckDB() (*sql.DB, error) {
	// Connecting to the database
	db, err := sql.Open("postgres", "host=db port=5432 user=postgres password=mysecretpassword dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Println("error connecting to the database: ", err)
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Checking if the database exists
	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Println("the database does not exist: ", err)
		return nil, fmt.Errorf("the database does not exist: %w", err)
	} else {
		log.Println("the database exists")
	}

	// Checking if the table exists
	_, err = db.Exec("SELECT 1 FROM fio_data LIMIT 1")
	if err != nil {
		log.Println("the table does not exist: ", err)
		return nil, fmt.Errorf("the table does not exist: %w", err)
	} else {
		log.Println("the table exists")
	}

	return db, nil
}
