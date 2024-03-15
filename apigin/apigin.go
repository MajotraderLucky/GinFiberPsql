package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/MajotraderLucky/Utils/logger"
	_ "github.com/lib/pq"
)

// Function to insert test data into the fio_data table
func InsertTestData(db *sql.DB) error {
	// SQL statement for inserting data
	query :=
		`INSERT INTO fio_data (name, surname, patronymic, age, gender, nationality) VALUES
($1, $2, $3, $4, $5, $6),
($7, $8, $9, $10, $11, $12),
($13, $14, $15, $16, $17, $18)`

	// Executing the query
	_, err := db.Exec(query,
		"Ivan", "Ivanov", "Ivanovich", 30, "male", "Russian",
		"Maria", "Petrova", "Sergeevna", 25, "female", "Russian",
		"John", "Doe", "Michael", 40, "male", "American",
	)
	if err != nil {
		log.Printf("Error inserting test data: %v", err)
		return fmt.Errorf("error inserting test data: %w", err)
	}

	log.Println("Test data inserted successfully")
	return nil
}

func main() {
	time.Sleep(time.Second * 10)

	logger := logger.Logger{}
	err := logger.CreateLogsDir()
	if err != nil {
		log.Fatal(err)
	}
	err = logger.OpenLogFile()
	if err != nil {
		log.Fatal(err)
	}
	logger.SetLogger()
	logger.LogLine()
	log.Println("Gin started...")
	logger.LogLine()

	// Database connection string
	connStr := "host=db port=5432 user=postgres password=mysecretpassword dbname=mydatabase sslmode=disable"

	// Opening connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Inserting test data
	err = InsertTestData(db)
	if err != nil {
		log.Fatal(err)
	}
}
