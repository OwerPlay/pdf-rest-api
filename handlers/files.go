package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"pdf-rest-api/database"

	"github.com/gorilla/mux"
)

func UploadFile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userId := mux.Vars(r)["userId"]
	if _, err := database.GetUser(db, userId); err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	// Read file from form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error": "File is required"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Enforce file size limit
	const maxFileSize = 10 << 20
	if handler.Size > maxFileSize {
		http.Error(w, `{"error": "File size exceeds 10MB limit"}`, http.StatusRequestEntityTooLarge)
		return
	}

	// Read first 512 bytes to determine MIME type
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		http.Error(w, `{"error": "Failed to read file"}`, http.StatusInternalServerError)
		return
	}
	file.Seek(0, io.SeekStart) // Reset file pointer

	// Validate MIME type
	if http.DetectContentType(buffer) != "application/pdf" {
		http.Error(w, `{"error": "Only PDF files are allowed"}`, http.StatusBadRequest)
		return
	}

	// Read the entire file into memory
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, `{"error": "Failed to read file"}`, http.StatusInternalServerError)
		return
	}

	// Store the file in the database
	fileID, err := database.InsertFile(db, userId, fileData, handler.Filename)
	if err != nil {
		http.Error(w, `{"error": "Database error while storing file"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"file_id":  fileID,
		"filename": handler.Filename,
		"message":  "File uploaded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteFile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract user ID and file ID from request
	userID := mux.Vars(r)["userId"]
	fileIDStr := mux.Vars(r)["fileId"]

	// Convert fileID to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Invalid file ID"}`, http.StatusBadRequest)
		return
	}

	// Attempt to delete the file
	err = database.DeleteFile(db, userID, fileID)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Failed to delete file"}`, http.StatusInternalServerError)
		return
	}

	// Send JSON response
	response := map[string]interface{}{
		"file_id": fileID,
		"message": "File deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ImportFile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract file ID from request parameters
	fileIDStr := mux.Vars(r)["fileId"]
	userID := mux.Vars(r)["userId"]

	// Check if the user exists
	if _, err := database.GetUser(db, userID); err != nil {
		http.Error(w, `{"status":"error","message":"User not found"}`, http.StatusNotFound)
		return
	}

	// TODO Check if the file exists
	// TODO Check if user has access to the file

	// Convert fileID to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Invalid file ID"}`, http.StatusBadRequest)
		return
	}

	// Attempt to import the file
	err = database.ImportFile(db, fileID)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Failed to import file"}`, http.StatusInternalServerError)
		return
	}

	// Send JSON response
	response := map[string]interface{}{
		"file_id": fileID,
		"message": "File imported successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
