package database

import (
	"database/sql"
	"errors"
	"log"
)

func CreateUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS USER (
		id INT AUTO_INCREMENT PRIMARY KEY
	);`
	_, err := db.Exec(query)
	if err != nil {
		return errors.New("Error creating USER table: " + err.Error())
	}

	log.Println("User table created successfully!")
	return nil
}

func CreateUser(db *sql.DB) (int, error) {
	result, err := db.Exec("INSERT INTO USER () VALUES ()")
	if err != nil {
		return 0, errors.New("Error inserting user: " + err.Error())
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("Error retrieving user ID: " + err.Error())
	}

	return int(userID), nil
}

func GetUser(db *sql.DB, userID string) (*string, error) {
	var UserID *string
	err := db.QueryRow("SELECT id FROM USER WHERE id = ?", userID).Scan(&UserID)
	if err != nil {
		return nil, errors.New("Error retrieving user: " + err.Error())
	}

	return UserID, nil
}
