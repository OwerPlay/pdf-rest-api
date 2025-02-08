package database

import (
	"database/sql"
	"log"
)

type StatusEnum string

const (
	StatusInQueue  StatusEnum = "in_queue"
	StatusParsing  StatusEnum = "parsing"
	StatusError    StatusEnum = "error"
	StatusSuccess  StatusEnum = "success"
	StatusImported StatusEnum = "imported"
)

type Files struct {
	ID         int        `json:"id"`
	Status     StatusEnum `json:"status"`
	PDFFile    []byte     `json:"pdf_file"`
	ParsedFile []byte     `json:"parsed_file"`
}

func CreateFilesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS FILES (
		id INT AUTO_INCREMENT PRIMARY KEY,
		status ENUM('in_queue', 'parsing', 'error', 'success', 'imported') NOT NULL,
		pdf_file LONGBLOB NOT NULL,
		parsed_file LONGBLOB
	);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Files table created successfully!")
	return nil
}
