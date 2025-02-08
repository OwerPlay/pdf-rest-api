package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// InitializeDatabase sets up the database and calls table creation functions
func InitializeDatabase() (*sql.DB, error) {
	// Read database credentials from environment variables
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := "mariadb"
	dbPort := "3306"
	dbName := os.Getenv("MYSQL_DATABASE")

	// Connect to MariaDB (without specifying a database initially)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MariaDB: %v", err)
	}

	// Create Database if not exists
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return nil, fmt.Errorf("error creating database: %v", err)
	}
	db.Close() // Close initial connection

	// Connect to the created database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Call table creation functions
	err = CreateUserTable(db)
	if err != nil {
		return nil, fmt.Errorf("error creating user table: %v", err)
	}

	err = CreateFilesTable(db)
	if err != nil {
		return nil, fmt.Errorf("error creating files table: %v", err)
	}

	err = CreateUserFilesTable(db)
	if err != nil {
		return nil, fmt.Errorf("error creating user_files table: %v", err)
	}

	log.Println("Database and tables initialized successfully!")
	return db, nil
}
