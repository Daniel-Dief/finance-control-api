package main

import (
	"finance-control-api/database"
	"finance-control-api/routes"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
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

	echoAPI := echo.New()

	routes.BindRoutes(echoAPI)

	if err := echoAPI.Start(":1323"); err != nil {
		echoAPI.Logger.Error("Failed to start server", "error", err)
	}
}
