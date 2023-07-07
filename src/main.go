package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Event struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
}

func main() {
	http.HandleFunc("/events", eventsHandler)
	http.HandleFunc("/events/", getEventHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getDBConnection() (*sql.DB, error) {
	// Create a MySQL database connection string
	connectionString := "root:password@tcp(localhost:3306)/sample"

	// Open a connection to the database
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func sendResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Write(jsonData)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getEventsHandler(w, r)
	case http.MethodPost:
		createEventHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the events table
	rows, err := db.Query("SELECT * FROM events")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store the events
	events := []Event{}

	// Iterate over the rows and append events to the slice
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.EventDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, event)
	}

	// Check for any errors occurred during iteration
	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the events response
	sendResponse(w, http.StatusOK, events)
}

func getEventHandler(w http.ResponseWriter, r *http.Request) {
	db, err := getDBConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Extract the event ID from the request URL
	eventID := strings.TrimPrefix(r.URL.Path, "/events/")
	if eventID == "" {
		http.Error(w, "Event ID not provided", http.StatusBadRequest)
		return
	}

	// Query the events table for the specific event
	row := db.QueryRow("SELECT * FROM events WHERE id = ?", eventID)

	// Create an Event struct to store the event data
	var event Event

	// Scan the result into the event struct
	err = row.Scan(&event.ID, &event.Name, &event.Description, &event.EventDate)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Send the event response
	sendResponse(w, http.StatusOK, event)
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into an Event struct
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a new UUID for the event ID
	event.ID = uuid.New().String()

	db, err := getDBConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert the new event into the events table
	_, err = db.Exec("INSERT INTO events (id, name, description, event_date) VALUES (?, ?, ?, ?)",
		event.ID, event.Name, event.Description, event.EventDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the success response
	sendResponse(w, http.StatusCreated, event)
}
