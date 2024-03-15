package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MajotraderLucky/Utils/logger"
	"github.com/gin-gonic/gin"
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

// Function to add data through the API
func addDataHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Structure for parsing the request body
		var data struct {
			Name        string `json:"name"`
			Surname     string `json:"surname"`
			Patronymic  string `json:"patronymic"`
			Age         int    `json:"age"`
			Gender      string `json:"gender"`
			Nationality string `json:"nationality"`
		}

		// Parsing the JSON body of the request
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// SQL query to add data
		query := `INSERT INTO fio_data (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := db.Exec(query, data.Name, data.Surname, data.Patronymic, data.Age, data.Gender, data.Nationality)
		if err != nil {
			log.Printf("Error inserting data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting data"})
			return
		}

		// Response about successful data addition
		c.JSON(http.StatusOK, gin.H{"message": "Data added successfully"})
	}
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

	// Creating a Gin router
	router := gin.Default()

	// Registering the handler for adding data
	router.POST("/add_data", addDataHandler(db))

	// Starting the server
	log.Fatal(router.Run(":8085"))
}
