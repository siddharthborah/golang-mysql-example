package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Create a MySQL database connection string
	connectionString := "root:password@tcp(localhost:3306)/sample"

	// Open a connection to the database
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Query the events table
	rows, err := db.Query("SELECT * FROM events")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the rows and print the values
	for rows.Next() {
		var id string
		var name string
		var description string
		var eventDate string
		err := rows.Scan(&id, &name, &description, &eventDate)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s, Name: %s, Description: %s, Date: %s\n", id, name, description, eventDate)
	}

	// Check for any errors occurred during iteration
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
