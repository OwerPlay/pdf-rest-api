package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"pdf-rest-api/config"
	"pdf-rest-api/database"
	"pdf-rest-api/handlers"
)

func main() {
	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/user/{userId}/files", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserFiles(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/queue", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetQueue(w, r, db)
	}).Methods("GET")

	r.HandleFunc("/user/create", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/user/{userId}/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadFile(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/file/{fileId}/parsed", func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadParsed(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/user/{userId}/file/{fileId}/import", func(w http.ResponseWriter, r *http.Request) {
		handlers.ImportFile(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/user/{userId}/file/{fileId}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteFile(w, r, db)
	}).Methods("DELETE")

	port := config.GetEnv("APP_PORT", "8080")
	log.Println("Server running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
