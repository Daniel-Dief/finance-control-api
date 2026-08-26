package main

import (
	"finance-control-api/database"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	fmt.Println("Environment variables loaded successfully.")

	ping, err := database.PingDatabase()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err.Error())
	}
	fmt.Println(ping)
}
