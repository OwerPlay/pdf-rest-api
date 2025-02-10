package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pdf-rest-api/database"

	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Create user in database
	userID, err := database.CreateUser(db)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"user_id": userID,
		"message": "User created successfully",
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetUserFiles(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract user ID from request parameters
	userID := mux.Vars(r)["userId"]

	// Check if the user exists
	if _, err := database.GetUser(db, userID); err != nil {
		http.Error(w, `{"status":"error","message":"User not found"}`, http.StatusNotFound)
		return
	}

	// Retrieve files for the user
	userFiles, err := database.GetUserFiles(db, userID)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Failed to retrieve user files"}`, http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userFiles)
}
