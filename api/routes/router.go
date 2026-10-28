package routes

import (
	"finance-control-api/database"
	"finance-control-api/middlewares"
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

	// Auth routes (public - no JWT required)
	authGroup := echoAPI.Group("/auth")
	Register(authGroup)
	Login(authGroup)

	// Create route groups for areas, budgets, categories, transactions, and reports
	areaGroup := echoAPI.Group("/areas")
	budgetGroup := echoAPI.Group("/budgets")
	categoryGroup := echoAPI.Group("/categories")
	transactionGroup := echoAPI.Group("/transactions")
	reportGroup := echoAPI.Group("/reports")

	// Apply JWT middleware to resource groups
	areaGroup.Use(middlewares.JWTMiddleware())
	budgetGroup.Use(middlewares.JWTMiddleware())
	categoryGroup.Use(middlewares.JWTMiddleware())
	transactionGroup.Use(middlewares.JWTMiddleware())
	reportGroup.Use(middlewares.JWTMiddleware())

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

	// Bind report routes
	GetKPIs(reportGroup)
	GetLineChart(reportGroup)
	GetRadarChart(reportGroup)
	GetBudgetReport(reportGroup)
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
