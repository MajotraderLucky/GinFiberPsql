package main

import (
	"ginfiberpsql/checkdb"
	"log"
	"time"

	"github.com/MajotraderLucky/Utils/logger"
	_ "github.com/lib/pq"
)

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

	// Close the database connection
	defer db.Close()

	logger.CleanLogCountLines(50)

}
