package database

import (
	"database/sql"
	"log"
)

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

func CreateUser(db *sql.DB) (int, error) {
	result, err := db.Exec("INSERT INTO USER () VALUES ()")
	if err != nil {
		log.Println("Error inserting user:", err)
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error retrieving user ID:", err)
		return 0, err
	}

	return int(userID), nil
}

func GetUser(db *sql.DB, userID string) (*string, error) {
	var UserID *string
	err := db.QueryRow("SELECT id FROM USER WHERE id = ?", userID).Scan(&UserID)
	if err != nil {
		log.Println("Error retrieving user:", err)
		return nil, err
	}

	return UserID, nil
}
