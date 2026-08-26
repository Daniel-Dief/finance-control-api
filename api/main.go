package main

import (
	"finance-control-api/database"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	log.Println("Environment variables loaded successfully.")

	err = database.InitPool()
	if err != nil {
		panic("Error initializing database pool: " + err.Error())
	}
	log.Println("Database pool initialized successfully.")

	err = database.PingDatabase()
	if err != nil {
		panic("Error pinging the database: " + err.Error())
	}
	log.Println("Database connection successful.")
}
