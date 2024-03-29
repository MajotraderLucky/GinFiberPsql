package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MajotraderLucky/Utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// Handler for querying a person by name and surname
func queryPersonHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Getting name and surname from the request query parameters
		name := c.Query("name")
		surname := c.Query("surname")

		// Checking if both name and surname are provided
		if name == "" || surname == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Both name and surname must be provided"})
			return
		}

		// Using the queryPersonByNameAndSurname function to query the database
		id, resultName, resultSurname, patronymic, age, gender, nationality, err := queryPersonByNameAndSurname(db, name, surname)
		if err != nil {
			// Handling the case when the person is not found or there is an error querying the database
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
			} else {
				log.Printf("Error querying person by name and surname: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying person"})
			}
			return
		}

		// Returning the found person's data
		c.JSON(http.StatusOK, gin.H{
			"id":          id,
			"name":        resultName,
			"surname":     resultSurname,
			"patronymic":  patronymic,
			"age":         age,
			"gender":      gender,
			"nationality": nationality,
		})
	}
}

var jwtKey = []byte("8GoPUxkoCEeKaEG381hL6p9RAfwgCaiDJhrwy+/k8Og=") // Replace with your key

// JWTAuthMiddleware checks for the presence and validity of a JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token signature algorithm is expected
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtKey, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
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
	router.POST("/add_data", JWTAuthMiddleware(), addDataHandler(db))
	router.GET("/query_person", JWTAuthMiddleware(), queryPersonHandler(db))

	// Starting the server
	log.Fatal(router.Run(":8085"))
}
