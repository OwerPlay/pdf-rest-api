package database

import (
	"database/sql"
	"errors"
	"log"
)

func CreateQueueTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS QUEUE (
		id INT AUTO_INCREMENT PRIMARY KEY,
		file_id INT NOT NULL,
		pdf_file LONGBLOB NOT NULL,
		FOREIGN KEY (file_id) REFERENCES FILES(id) ON DELETE CASCADE
	);`
	_, err := db.Exec(query)
	if err != nil {
		return errors.New("Error creating QUEUE table: " + err.Error())
	}

	log.Println("QUEUE table created successfully!")
	return nil
}

func AddFileToQueue(db *sql.DB, fileID int, fileData []byte) error {
	_, err := db.Exec("INSERT INTO QUEUE (file_id, pdf_file) VALUES (?, ?)", fileID, fileData)
	if err != nil {
		return errors.New("Error adding file to queue: " + err.Error())
	}
	log.Println("File added to queue:", fileID)
	return nil
}

func GetNextFileFromQueue(db *sql.DB) (int, []byte, error) {
	var fileID int
	var fileData []byte
	query := "SELECT file_id, pdf_file FROM QUEUE ORDER BY id ASC LIMIT 1"

	err := db.QueryRow(query).Scan(&fileID, &fileData)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Queue is empty.")
			return 0, nil, nil
		}
		return 0, nil, errors.New("Error fetching next file from queue: " + err.Error())
	}

	// Remove file from queue after fetching
	_, err = db.Exec("DELETE FROM QUEUE WHERE file_id = ? LIMIT 1", fileID)
	if err != nil {
		return 0, nil, errors.New("Error removing file from queue: " + err.Error())
	}

	// Update file status to 'parsing'
	err = UpdateFileStatus(db, fileID, StatusParsing)
	if err != nil {
		return 0, nil, errors.New("Error updating file status: " + err.Error())
	}

	log.Println("File removed from queue:", fileID)
	return fileID, fileData, nil
}
