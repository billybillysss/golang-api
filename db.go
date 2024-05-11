package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

// connDb establishes a connection to the PostgreSQL database
func connDb() {
	// Construct the PostgreSQL connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Log the connection string
	fmt.Println("DEBUG: Connection String:", psqlInfo)

	// Attempt to open a connection to the database
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		// If there was an error opening the connection, log it and panic
		fmt.Println("ERROR: Failed to connect to database:", err)
		panic(err)
	}

	// Log the successful connection
	fmt.Println("DEBUG: Establishing connection to database...")

	// Attempt to ping the database
	err = db.Ping()
	if err != nil {
		// If there was an error pinging the database, log it and panic
		fmt.Println("ERROR: Failed to connect to database:", err)
		panic(err)
	}

	// Log the successful ping
	fmt.Println("DEBUG: Established a successful connection!")
}
