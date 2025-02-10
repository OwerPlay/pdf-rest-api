package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitializeDatabase() (*sql.DB, error) {
	// Read database credentials from environment variables
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbHost := "mariadb"
	dbPort := "3306"
	dbName := os.Getenv("MYSQL_DATABASE")

	// Connect to MariaDB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.New("error connecting to database")
	}

	// Create Database if not exists
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return nil, errors.New("error creating database: " + err.Error())
	}
	db.Close()

	// Connect to the created database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.New("error connecting to database: " + err.Error())
	}

	// Call table creation functions
	err = CreateUserTable(db)
	if err != nil {
		return nil, errors.New("error creating user table: " + err.Error())
	}

	err = CreateFilesTable(db)
	if err != nil {
		return nil, errors.New("error creating files table: " + err.Error())
	}

	err = CreateUserFilesTable(db)
	if err != nil {
		return nil, errors.New("error creating user_files table: " + err.Error())
	}

	err = CreateQueueTable(db)
	if err != nil {
		return nil, errors.New("error creating queue table: " + err.Error())
	}

	log.Println("Database and tables initialized successfully!")
	return db, nil
}
