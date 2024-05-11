package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
}

// getMember retrieves a member or multiple members from the database
//
// If no member ID is provided, it retrieves all members
// If one or more member IDs are provided, it retrieves those members
//
// Parameters:
//
//	id ...int - The ID of the member(s) to retrieve
//
// Returns:
//
//	[]Member - The member(s) retrieved from the database
//	error - Any error that may have occurred
func getMember(id ...int) ([]Member, error) {
	var members []Member
	var member Member

	// Create the SELECT SQL statement
	sqlQuery := selectSql(member, "members", id...)

	// Log the SQL query being executed
	log.Println("Executing SQL query:", sqlQuery)

	// Execute the SQL query
	rows, err := db.Query(sqlQuery)
	if err != nil {
		// If there is an error executing the query, return the error
		log.Println("Error executing query:", err.Error())
		return members, err
	}
	defer rows.Close() // Close the rows result set when finished

	// Convert the rows result set to JSON
	// err = convertToJson(rows, member.Fields())
	for rows.Next() {
		if err := rows.Scan(member.Fields()...); err != nil {
			fmt.Println(err.Error())
			return members, err
		}
		members = append(members, member)
	}
	// Return the members retrieved from the database
	return members, err
}

// getMemberHandle handles GET requests to /members/{member_id}
// This function gets a single member from the database
// based on the member_id in the URL
func getMemberHandle(w http.ResponseWriter, r *http.Request) {
	// Set the content type of the response to JSON
	w.Header().Set("Content-Type", "application/json")

	// Declare a variable to store the member struct
	var members []Member
	// Declare a variable to store any error that may occur
	var err error

	// Get the member ID from the URL
	params := mux.Vars(r)
	memberId, err := strconv.Atoi(params["member_id"])
	if err != nil {
		// If there is an error converting the member ID to an int, return a failure message
		log.Println("Error converting member_id to int:", err.Error())
		response := Response{Message: "Failed to get member!"}
		json.NewEncoder(w).Encode(response)
		panic(err)
	}

	// Log the member ID being requested
	log.Println("Getting member with ID:", memberId)

	// Get the member from the database
	members, err = getMember(memberId)
	if err != nil {
		// If there is an error getting the member, return a failure message
		log.Println("Error getting member:", err.Error())
		response := Response{Message: "Failed to get member!"}
		json.NewEncoder(w).Encode(response)
		panic(err)
	}

	// If there is no error, return the member
	log.Println("Returning member:", members[0])
	json.NewEncoder(w).Encode(members[0])
}

// getMembersHandle handles GET requests to /members
// This function gets all members from the database
func getMembersHandle(w http.ResponseWriter, r *http.Request) {
	// Set the content type of the response to JSON
	w.Header().Set("Content-Type", "application/json")

	// Declare a variable to store the member structs
	var members []Member
	// Declare a variable to store any error that may occur
	var err error

	// Log the SQL statement being executed
	log.Println("Getting all members")

	// Get all members from the database
	members, err = getMember()
	if err != nil {
		// If there is an error, return a failure message
		log.Println("Error getting members:", err.Error())
		response := Response{Message: "Failed to get members!"}
		json.NewEncoder(w).Encode(response)
		panic(err)
	}

	// If there is no error, return the members
	log.Println("Returning members:", members)
	json.NewEncoder(w).Encode(members)
}

// CreateMemberHandle handles POST requests to /members
// This function creates a new member in the database
func createMemberHandle(w http.ResponseWriter, r *http.Request) {
	// Set the content type of the response to JSON
	w.Header().Set("Content-Type", "application/json")

	// Declare a variable to store the member struct
	var member Member

	// Decode the JSON body of the request into the member struct
	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		// If there is an error decoding the JSON body, return a failure message
		log.Println("Error decoding JSON body:", err.Error())
		response := Response{Message: "Failed to decode JSON body!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Marshal the member struct to JSON
	data, err := json.Marshal(member)
	if err != nil {
		// If there is an error marshalling the member, return a failure message
		log.Println("Error marshalling member:", err.Error())
		panic(err)
	}

	// Create an INSERT SQL statement to insert the member
	sqlScript := updateOrInsertSql("members", string(data), "insert")

	// Log the SQL statement being executed
	log.Println("Executing SQL:", sqlScript)

	// Execute the SQL statement
	_, err = db.Query(sqlScript)
	if err != nil {
		// If there is an error executing the SQL statement, return a failure message
		log.Println("Error inserting member:", err.Error())
		response := Response{Message: "Failed to insert!"}
		json.NewEncoder(w).Encode(response)
		panic(err)
	}

	// If there is no error, return a success message
	log.Println("Inserted member successfully!")
	response := Response{Message: "Success!"}
	json.NewEncoder(w).Encode(response)
}

// UpdateMemberHandle handles PUT requests to /members/{member_id}
// This function updates a member in the database
func updateMemberHandle(w http.ResponseWriter, r *http.Request) {
	// Set the content type of the response to JSON
	w.Header().Set("Content-Type", "application/json")

	// Declare a variable to store the member struct
	var member Member

	// Get the member ID from the URL
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["member_id"])
	if err != nil {
		// If there is an error converting the member ID to an int, return a failure message
		log.Println("Error converting member_id to int:", err.Error())
		response := Response{Message: "Failed! ID mismatch"}
		json.NewEncoder(w).Encode(response)
		panic(err)
	}
	log.Println("Updating member with ID:", id)

	// Decode the JSON body of the request into the member struct
	json.NewDecoder(r.Body).Decode(&member)
	log.Println("Member object:", member)

	// Marshal the member struct to JSON
	data, err := json.Marshal(member)
	if err != nil {
		// If there is an error marshalling the member, return a failure message
		log.Println("Error marshalling member:", err.Error())
		log.Fatal(err.Error())
	}

	// Check if the member ID in the URL matches the member ID in the JSON
	if id != member.MemberID {
		// If the IDs don't match, return a failure message
		log.Println("ID mismatch:", id, "!=", member.MemberID)
		response := Response{Message: "Failed! ID mismatch"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create an UPDATE SQL statement to update the member
	sqlScript := updateOrInsertSql("members", string(data), "update")
	log.Println("Executing SQL:", sqlScript)

	// Execute the SQL statement
	_, err = db.Query(sqlScript)
	if err != nil {
		// If there is an error executing the SQL statement, return a failure message
		log.Println("Error updating member:", err.Error())
		panic(err)
	}

	// If there is no error, return a success message
	log.Println("Updated member successfully!")
	response := Response{Message: "Success!"}
	json.NewEncoder(w).Encode(response)

	return
}

// DeleteMemberHandle handles DELETE requests to /members/{member_id}
// This function deletes a member from the database
func deleteMemberHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Create a DELETE SQL statement to delete a member
	sqlScript := deleteSql(Member{}, "members", 1)

	// Log the SQL statement being executed
	log.Println("Executing SQL:", sqlScript)

	// Execute the SQL statement
	_, err := db.Query(sqlScript)
	if err != nil {
		// If there is an error, return a failure message
		log.Println("Error deleting member:", err.Error())
		response := Response{Message: "Failed to delete!"}
		json.NewEncoder(w).Encode(response)
		panic(err)
	} else {
		// If there is no error, return a success message
		log.Println("Deleted member successfully!")
		response := Response{Message: "Success!"}
		json.NewEncoder(w).Encode(response)
	}
}
