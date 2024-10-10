package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func runMigration(db *sql.DB) error {
	// Read the SQL file
	data, err := ioutil.ReadFile("create_tables.sql")
	if err != nil {
		return err
	}

	// Execute the SQL commands
	_, err = db.Exec(string(data))
	return err
}

// openDB initializes a connection to the database using the provided DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB establishes a connection to the database using environment variables.
func (app *application) connectToDB() (*sql.DB, error) {
	// Retrieve database connection details from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
		host, port, user, password, dbname)

	connection, err := openDB(dsn)
	if err != nil {
		return nil, err
	}

	if err := runMigration(connection); err != nil {
		return nil, err
	}

	log.Println("Connected to the database!")
	return connection, nil
}
