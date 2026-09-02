package main

import (
	"finance-control-api/database"
	_ "finance-control-api/docs"
	"finance-control-api/routes"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

// @title Finance Control API
// @version 1.0
// @description API para gerenciamento financeiro pessoal — áreas, categorias, orçamentos e transações.
// @host localhost:1323
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}
	log.Println("Environment variables loaded successfully.")

	err = database.InitPool()
	if err != nil {
		log.Fatal("Error initializing database pool: " + err.Error())
	}
	log.Println("Database pool initialized successfully.")

	err = database.PingDatabase()
	if err != nil {
		log.Fatal("Error pinging the database: " + err.Error())
	}
	log.Println("Database connection successful.")

	echoAPI := echo.New()

	routes.BindRoutes(echoAPI)

	if err := echoAPI.Start(":1323"); err != nil {
		echoAPI.Logger.Error("Failed to start server", "error", err)
	}
}
