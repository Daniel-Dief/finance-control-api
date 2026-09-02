package routes

import (
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger/v2"
)

func BindRoutes(echoAPI *echo.Echo) {
	echoAPI.GET("/swagger/*", echoSwagger.WrapHandler)

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
