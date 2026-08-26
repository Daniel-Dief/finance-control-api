package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var Pool *sql.DB

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
func PingDatabase() error {
	if Pool == nil {
		return fmt.Errorf("Database pool is not initialized")
	}

	err := Pool.Ping()
	if err != nil {
		return fmt.Errorf("Database connection failed: %v", err)
	}

	return nil
}

// InitPool initializes the database connection pool.
func InitPool() error {
	strConn := StringConn()

	db, err := sql.Open("postgres", strConn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	Pool = db

	return nil
}
