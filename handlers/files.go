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
		http.Error(w, `{"error": "User not found", "details": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	// Read file from form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error": "File is required", "details": "`+err.Error()+`"}`, http.StatusBadRequest)
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
		http.Error(w, `{"error": "Failed to read file", "details": "`+err.Error()+`"}`, http.StatusInternalServerError)
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
		http.Error(w, `{"error": "Failed to read file", "details": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Store the file in the database
	fileID, err := database.InsertFile(db, userId, fileData, handler.Filename)
	if err != nil {
		http.Error(w, `{"error": "Database error while storing file", "details": "`+err.Error()+`"}`, http.StatusInternalServerError)
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
		http.Error(w, `{"status":"error","message":"Invalid file ID", "details": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Attempt to delete the file
	err = database.DeleteFile(db, userID, fileID)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Failed to delete file", "details": "`+err.Error()+`"}`, http.StatusInternalServerError)
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
		http.Error(w, `{"status":"error","message":"User not found", "details": "`+err.Error()+`"}`, http.StatusNotFound)
		return
	}

	// Convert fileID to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Invalid file ID", "details": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Check if user with this file exists
	exists, err := database.UserFileExists(db, userID, fileID)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Failed to check user-file existence", "details": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, `{"status":"error","message":"User with this file does not exist"}`, http.StatusBadRequest)
		return
	}

	// Attempt to import the file
	err = database.ImportFile(db, fileID)
	if err != nil {
		http.Error(w, `{"status":"error","message":"Failed to import file", "details": "`+err.Error()+`"}`, http.StatusInternalServerError)
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
