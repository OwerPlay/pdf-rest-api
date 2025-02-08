package database

import (
	"database/sql"
	"log"
)

type User struct {
	ID int `json:"id"`
}

func CreateUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS USER (
		id INT AUTO_INCREMENT PRIMARY KEY
	);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("User table created successfully!")
	return nil
}
