package main

import (
	"log"
	"net/http"

	"pdf-rest-api/config"
	"pdf-rest-api/database"
	"pdf-rest-api/handlers"
)

func main() {
	// Initialize Database
	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// Read the port from environment variables, default to 8080
	port := config.GetEnv("APP_PORT", "8080")

	// Register handlers
	http.HandleFunc("/test", handlers.TestHandler)

	log.Println("Server running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
