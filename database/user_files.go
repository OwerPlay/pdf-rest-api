package database

import (
	"database/sql"
	"log"
)

type UserFiles struct {
	UserID     int    `json:"user_id"`
	FileID     int    `json:"file_id"`
	Filename   string `json:"filename"`
	UploadDate string `json:"upload_date"`
}

// CreateUserFilesTable creates the USER_FILES table
func CreateUserFilesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS USER_FILES (
		user_id INT NOT NULL,
		file_id INT NOT NULL,
		filename VARCHAR(255) NOT NULL,
		upload_date DATE NOT NULL,
		PRIMARY KEY (user_id, file_id),
		FOREIGN KEY (user_id) REFERENCES USER(id) ON DELETE CASCADE,
		FOREIGN KEY (file_id) REFERENCES FILES(id) ON DELETE CASCADE
	);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("User_Files table created successfully!")
	return nil
}
