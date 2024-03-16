package apifunc

import (
	"database/sql"
	"fmt"
	"log"
)

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
