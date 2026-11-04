package main

import (
	"finance-control-api/database"
	_ "finance-control-api/docs"
	"finance-control-api/helpers"
	"finance-control-api/middlewares"
	"finance-control-api/routes"
	"log"

	"github.com/labstack/echo/v5"
)

// @title Finance Control API
// @version 1.0
// @description API para gerenciamento financeiro pessoal — áreas, categorias, orçamentos e transações.
// @host localhost:1323
// @BasePath /
func main() {
	if err := helpers.CheckEnvs(); err != nil {
		log.Fatal("Environment check failed: " + err.Error())
	} else {
		log.Println("Environment variables loaded successfully.")
	}

	if err := database.InitPool(); err != nil {
		log.Fatal("Error initializing database pool: " + err.Error())
	} else {
		log.Println("Database pool initialized successfully.")
	}

	if err := database.PingDatabase(); err != nil {
		log.Fatal("Error pinging the database: " + err.Error())
	} else {
		log.Println("Database connection successful.")
	}

	echoAPI := echo.New()

	middlewares.BindMiddlewares(echoAPI)
	routes.BindRoutes(echoAPI)

	if err := echoAPI.Start(":1323"); err != nil {
		echoAPI.Logger.Error("Failed to start server", "error", err)
	}
}
