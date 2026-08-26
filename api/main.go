package main

import (
	"finance-control-api/database"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}
	fmt.Println("Environment variables loaded successfully.")

	ping, err := database.PingDatabase()
	if err != nil {
		panic("Error pinging the database: " + err.Error())
	}
	fmt.Println(ping)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Failed to start server: ", "error", err)
	}
}
