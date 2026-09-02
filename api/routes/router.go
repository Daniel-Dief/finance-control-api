package routes

import (
	"finance-control-api/database"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func BindRoutes(echoAPI *echo.Echo) {
	if os.Getenv("ENV") == "development" {
		echoAPI.GET("/swagger/*", echoSwagger.WrapHandler)
	}
	// Bind health check route
	HealthCheck(echoAPI)

	// Create route groups for areas, budgets, categories, and transactions
	areaGroup := echoAPI.Group("/areas")
	budgetGroup := echoAPI.Group("/budgets")
	categoryGroup := echoAPI.Group("/categories")
	transactionGroup := echoAPI.Group("/transactions")

	// Bind area routes
	ListAreas(areaGroup)
	GetAreaByID(areaGroup)
	CreateArea(areaGroup)
	UpdateArea(areaGroup)
	DeleteArea(areaGroup)

	// Bind category routes
	ListCategories(categoryGroup)
	GetCategoryByID(categoryGroup)
	CreateCategory(categoryGroup)
	UpdateCategory(categoryGroup)
	DeleteCategory(categoryGroup)

	// Bind budget routes
	ListBudgets(budgetGroup)
	GetBudgetByID(budgetGroup)
	CreateBudget(budgetGroup)
	UpdateBudget(budgetGroup)
	DeleteBudget(budgetGroup)

	// Bind transaction routes
	ListTransactions(transactionGroup)
	GetTransactionByID(transactionGroup)
	CreateTransaction(transactionGroup)
	UpdateTransaction(transactionGroup)
	DeleteTransaction(transactionGroup)
}

func HealthCheck(echoAPI *echo.Echo) {
	echoAPI.GET("/health", func(c *echo.Context) error {
		err := database.PingDatabase()

		if err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "System down"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
}
