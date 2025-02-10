package database

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
)

type StatusEnum string

const (
	StatusInQueue  StatusEnum = "in_queue"
	StatusParsing  StatusEnum = "parsing"
	StatusError    StatusEnum = "error"
	StatusSuccess  StatusEnum = "success"
	StatusImported StatusEnum = "imported"
)

func CreateFilesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS FILES (
		id INT AUTO_INCREMENT PRIMARY KEY,
		status ENUM('in_queue', 'parsing', 'error', 'success', 'imported') NOT NULL,
		parsed_file LONGBLOB,
		file_hash VARCHAR(64) UNIQUE NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		return errors.New("Error creating files table: " + err.Error())
	}

	return nil
}

func InsertFile(db *sql.DB, userId string, fileData []byte, filename string) (int, error) {
	fileHash := generateFileHash(fileData)

	// Check if file already exists
	var existingFileID int
	err := db.QueryRow("SELECT id FROM FILES WHERE file_hash = ?", fileHash).Scan(&existingFileID)

	if err == nil {
		// Check if user already has access to this file
		err = InsertUserFile(db, userId, existingFileID, filename)
		if err != nil {
			return 0, errors.New("Error inserting user file: " + err.Error())
		}
		return existingFileID, nil
	} else if err != sql.ErrNoRows {
		return 0, errors.New("Error checking file existence: " + err.Error())
	}

	// Insert new file
	result, err := db.Exec("INSERT INTO FILES (status, file_hash) VALUES (?, ?)", StatusInQueue, fileHash)
	if err != nil {
		return 0, errors.New("Error inserting file: " + err.Error())
	}

	// Retrieve the file ID
	fileID, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("Error retrieving file ID: " + err.Error())
	}

	// Store the user-file relationship
	err = InsertUserFile(db, userId, int(fileID), filename)
	if err != nil {
		return 0, errors.New("Error inserting user file: " + err.Error())
	}

	// Add file to QUEUE with its actual file data
	err = AddFileToQueue(db, int(fileID), fileData)
	if err != nil {
		return 0, errors.New("Error adding file to queue: " + err.Error())
	}

	return int(fileID), nil
}

func UpdateFileStatus(db *sql.DB, fileID int, status StatusEnum) error {
	_, err := db.Exec("UPDATE FILES SET status = ? WHERE id = ?", status, fileID)
	if err != nil {
		return errors.New("Error updating file status: " + err.Error())
	}

	return nil
}

func UploadParsedFile(db *sql.DB, fileID int, parsedFile []byte) error {
	if parsedFile == nil || len(parsedFile) == 0 {
		_, err := db.Exec("UPDATE FILES SET status = ? WHERE id = ?", StatusError, fileID)
		if err != nil {
			return errors.New("Error uploading parsed file: " + err.Error())
		}
	} else {
		_, err := db.Exec("UPDATE FILES SET status = ?, parsed_file = ? WHERE id = ?", StatusSuccess, parsedFile, fileID)
		if err != nil {
			return errors.New("Error uploading parsed file: " + err.Error())
		}
	}

	return nil
}

func ImportFile(db *sql.DB, fileID int) error {
	// Check if file is in 'success' state
	var status StatusEnum
	err := db.QueryRow("SELECT status FROM FILES WHERE id = ?", fileID).Scan(&status)
	if err != nil {
		return errors.New("Error retrieving file status: " + err.Error())
	}

	if status != StatusSuccess {
		return errors.New("File is not in 'success' state. Cannot import.")
	}

	err = UpdateFileStatus(db, fileID, StatusImported)
	if err != nil {
		return errors.New("Error importing file: " + err.Error())
	}

	return nil
}

func generateFileHash(fileData []byte) string {
	hash := sha256.Sum256(fileData)
	return hex.EncodeToString(hash[:])
}
