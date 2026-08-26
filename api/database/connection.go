package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// StringConn constructs the PostgreSQL connection string using environment variables.
func StringConn() string {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",

		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	return psqlInfo
}

// PingDatabase pings the database to check if the connection is successful.
func PingDatabase() (string, error) {
	db, err := CreatePool()
	if err != nil {
		return "", err
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		return "", err
	}

	return "Database connection successful!", nil
}

// CreatePool creates a connection pool to the PostgreSQL database and returns the *sql.DB instance.
// dont forget to close the pool when you're done using it by calling db.Close().
func CreatePool() (*sql.DB, error) {
	strConn := StringConn()
	db, err := sql.Open("postgres", strConn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}
