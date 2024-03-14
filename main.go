package main

import (
	"database/sql"
	"ginfiberpsql/checkdb"
	"log"
	"time"

	"github.com/MajotraderLucky/Utils/logger"
	_ "github.com/lib/pq"
)

// Function to execute a query using a composite index
func queryPersonByNameAndSurname(db *sql.DB, name, surname string) (int, string, string, string, int, string, string, error) {
	// The query uses the composite index idx_fio_data_name_surname
	query := `SELECT id, name, surname, patronymic, age, gender, nationality FROM fio_data WHERE name = $1 AND surname = $2;`

	// Executing the query
	row := db.QueryRow(query, name, surname)

	// Reading the results
	var (
		id            int
		resultName    string
		resultSurname string
		patronymic    string
		age           int
		gender        string
		nationality   string
	)

	err := row.Scan(&id, &resultName, &resultSurname, &patronymic, &age, &gender, &nationality)
	if err != nil {
		return 0, "", "", "", 0, "", "", err
	}

	return id, resultName, resultSurname, patronymic, age, gender, nationality, nil
}

func main() {
	// Wait for other services to be ready
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

	log.Println("GinFiberPsql started...")
	logger.LogLine()

	// Connect to the database and check it
	db, err := checkdb.ConnectAndCheckDB()
	if err != nil {
		log.Fatal(err)
	}

	// Close the database connection at the end
	defer db.Close()

	// Filter values
	name := "Иван"
	surname := "Иванов"

	// Using the function to execute the query
	id, resultName, resultSurname, patronymic, age, gender, nationality, err := queryPersonByNameAndSurname(db, name, surname)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ID: %d, Name: %s, Surname: %s, Patronymic: %s, Age: %d, Gender: %s, Nationality: %s\n",
		id, resultName, resultSurname, patronymic, age, gender, nationality)

	logger.CleanLogCountLines(50)

}
