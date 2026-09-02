package tests

import (
	"finance-control-api/database"
	"finance-control-api/models"
	"finance-control-api/routes"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	err = database.InitPool()
	if err != nil {
		log.Fatal("Error initializing database pool: " + err.Error())
	}

	err = database.PingDatabase()
	if err != nil {
		log.Fatal("Error pinging the database: " + err.Error())
	}

	os.Exit(m.Run())
}

func setupEcho() *echo.Echo {
	e := echo.New()
	routes.BindRoutes(e)
	return e
}

func createTestArea(name string) (models.Area, error) {
	return database.CreateArea(name)
}

func createTestBudget(year, month, areaID, amount int) (models.Budget, error) {
	props := database.BudgetProps{
		Year:   &year,
		Month:  &month,
		AreaId: &areaID,
		Amount: &amount,
	}
	return database.CreateBudget(props)
}

func createTestCategory(name string) (models.Category, error) {
	return database.CreateCategory(name)
}

func createTestTransaction(date string, amount, categoryID, areaID int, transType string) (models.Transaction, error) {
	props := database.TransactionProps{
		Date:       &date,
		Amount:     &amount,
		CategoryID: &categoryID,
		AreaID:     &areaID,
		Type:       &transType,
	}
	return database.CreateTransaction(props)
}