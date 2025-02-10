package database

import (
	"database/sql"
	"log"
	"time"
)

type UserFileDetails struct {
	UploadDate string `json:"upload_date"`
	Filename   string `json:"filename"`
	Status     string `json:"status"`
}

// CreateUserFilesTable creates the USER_FILES table with TIMESTAMP support
func CreateUserFilesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS USER_FILES (
		user_id INT NOT NULL,
		file_id INT NOT NULL,
		filename VARCHAR(255) NOT NULL,
		upload_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, file_id),
		FOREIGN KEY (user_id) REFERENCES USER(id) ON DELETE CASCADE,
		FOREIGN KEY (file_id) REFERENCES FILES(id) ON DELETE CASCADE
	);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("User_Files table created successfully with TIMESTAMP support!")
	return nil
}

// InsertUserFile creates a user-file relationship if it doesn't already exist
func InsertUserFile(db *sql.DB, userId string, fileID int, filename string) error {
	// Check if the user-file relationship already exists
	exists, err := UserFileExists(db, userId, fileID)
	if err != nil {
		return err
	}
	if exists {
		log.Println("User already has access to this file.")
		return nil
	}

	// Insert new user-file link with current timestamp
	_, err = db.Exec("INSERT INTO USER_FILES (user_id, file_id, filename, upload_date) VALUES (?, ?, ?, ?)",
		userId, fileID, filename, time.Now().Format("2006-01-02 15:04:05"))

	if err != nil {
		log.Println("Error linking file to user:", err)
		return err
	}

	log.Println("File linked to user:", userId)
	return nil
}

// UserFileExists checks if a user already has access to a specific file
func UserFileExists(db *sql.DB, userId string, fileID int) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM USER_FILES WHERE user_id = ? AND file_id = ?", userId, fileID).Scan(&exists)
	if err != nil {
		log.Println("Error checking user-file existence:", err)
		return false, err
	}
	return exists > 0, nil
}

func GetUserFiles(db *sql.DB, userID string) ([]UserFileDetails, error) {
	query := `
		SELECT uf.upload_date, uf.filename, f.status
		FROM USER_FILES uf
		JOIN FILES f ON uf.file_id = f.id
		WHERE uf.user_id = ?
		ORDER BY uf.upload_date DESC;
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error retrieving user files:", err)
		return nil, err
	}
	defer rows.Close()

	var userFiles []UserFileDetails
	for rows.Next() {
		var file UserFileDetails
		err := rows.Scan(&file.UploadDate, &file.Filename, &file.Status)
		if err != nil {
			log.Println("Error scanning user file row:", err)
			return nil, err
		}
		userFiles = append(userFiles, file)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error iterating user files rows:", err)
		return nil, err
	}

	return userFiles, nil
}

func DeleteFile(db *sql.DB, userId string, fileID int) error {
	// Ensure the file is in "in_queue" state
	var status string
	err := db.QueryRow("SELECT status FROM FILES WHERE id = ?", fileID).Scan(&status)
	if err != nil {
		log.Println("Error retrieving file status:", err)
		return err
	}

	if status != "in_queue" {
		log.Println("File cannot be deleted as it is not in 'in_queue' state.")
		// TODO return correct error check this everywhere
		return sql.ErrNoRows
	}

	// Delete user-file relationship
	_, err = db.Exec("DELETE FROM USER_FILES WHERE user_id = ? AND file_id = ?", userId, fileID)
	if err != nil {
		log.Println("Error removing user-file link:", err)
		return err
	}
	log.Println("User-file link removed:", userId, fileID)

	// Check if the file is still linked to any users
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM USER_FILES WHERE file_id = ?", fileID).Scan(&count)
	if err != nil {
		log.Println("Error checking remaining user links for file:", err)
		return err
	}

	// If no users are linked to the file, delete it from FILES
	if count == 0 {
		_, err = db.Exec("DELETE FROM FILES WHERE id = ?", fileID)
		if err != nil {
			log.Println("Error deleting file from FILES:", err)
			return err
		}
		log.Println("File deleted from FILES:", fileID)
	}

	return nil
}
