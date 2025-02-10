package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pdf-rest-api/database"
	"strconv"

	"github.com/gorilla/mux"
)

type ParsedFilePayload struct {
	ParsedFile   []byte `json:"parsed_file"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

func GetQueue(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Retrieve the next file in the queue
	fileID, fileData, err := database.GetNextFileFromQueue(db)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Failed to retrieve next file", "details": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	if fileID == 0 {
		http.Error(w, `{"status":"error","message":"Queue is empty"}`, http.StatusNotFound)
		return
	}

	// Prepare response
	response := map[string]interface{}{
		"file_id":  fileID,
		"pdf_file": fileData, // Returning binary data in JSON isn't ideal; consider Base64 encoding
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UploadParsed(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract file ID from request parameters
	fileIDStr := mux.Vars(r)["fileId"]

	// Convert fileID to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Invalid file ID", "details": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Decode JSON payload
	var payload ParsedFilePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, `{"status":"error","message":"Invalid JSON payload", "details": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Set status to error if errorMessage is provided
	if payload.ErrorMessage != "" {
		payload.Status = string(database.StatusError)
	}

	// Store the parsed file in the database
	err = database.UploadParsedFile(db, fileID, payload.ParsedFile, payload.Status)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Database error while storing file", "details": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Send JSON response
	response := map[string]interface{}{
		"message": "Parsed information uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
