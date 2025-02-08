package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type TestResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request method"})
		return
	}

	// Prepare the test response
	response := TestResponse{
		Status:  "success",
		Message: "Test complete",
	}

	// Set response headers and return JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	log.Println("Processing completed successfully")
}
