package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// This is the entry point of the application. When we run the Go program,
// this is where it will start.

func main() {

	// Connect to the database
	connDb()
	defer db.Close() // Ensure the database connection is closed when the function exits

	// Create a new router
	r := mux.NewRouter()

	// Handle GET requests to the /members endpoint
	r.HandleFunc("/members", getMembersHandle).Methods("GET")
	// Handle GET requests to the /members/{member_id} endpoint
	r.HandleFunc("/members/{member_id:[0-9]+}", getMemberHandle).Methods("GET")
	// Handle POST requests to the /members endpoint
	r.HandleFunc("/members", createMemberHandle).Methods("POST")
	// Handle PUT requests to the /members/{member_id} endpoint
	r.HandleFunc("/members/{member_id:[0-9]+}", updateMemberHandle).Methods("PUT")
	// Handle DELETE requests to the /members/{member_id} endpoint
	r.HandleFunc("/members/{member_id:[0-9]+}", deleteMemberHandle).Methods("DELETE")

	// Start the server and log any errors
	log.Fatal(http.ListenAndServe(":8000", r))

}
