package routes

import (
	"finance-control-api/database"
	"finance-control-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// ListBudgets defines the route for listing budgets.
func ListBudgets(budgetGroup *echo.Group) {
	budgetGroup.GET("/list", func(c *echo.Context) error {
		props := database.BudgetProps{
			Year:   parseIntPtr(c.QueryParam("year")),
			Month:  parseIntPtr(c.QueryParam("month")),
			AreaId: parseIntPtr(c.QueryParam("area_id")),
		}

		budgets, err := database.ListBudgets(props)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if len(budgets) == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Nenhum orçamento encontrado"})
		}

		return c.JSON(http.StatusOK, budgets)
	})
}

// GetBudgetByID defines the route for retrieving a budget by its ID.
func GetBudgetByID(budgetGroup *echo.Group) {
	budgetGroup.GET("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		budget, err := database.GetBudgetByID(id)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if budget.ID == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Orçamento não encontrado"})
		}

		return c.JSON(http.StatusOK, budget)
	})
}

// CreateBudget defines the route for creating a new budget.
func CreateBudget(budgetGroup *echo.Group) {
	budgetGroup.POST("/create", func(c *echo.Context) error {
		var toCreateBudget models.Budget

		if err := c.Bind(&toCreateBudget); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		}

		if toCreateBudget.AreaID == 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O id da área é obrigatório"})
		}
		if toCreateBudget.Year == 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O ano é obrigatório"})
		}
		if toCreateBudget.Month == 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O mês é obrigatório"})
		}

		props := database.BudgetProps{
			Year:   &toCreateBudget.Year,
			Month:  &toCreateBudget.Month,
			AreaId: &toCreateBudget.AreaID,
			Amount: &toCreateBudget.Amount,
		}

		createdBudget, err := database.CreateBudget(props)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, createdBudget)
	})
}

// UpdateBudget defines the route for updating an existing budget.
func UpdateBudget(budgetGroup *echo.Group) {
	budgetGroup.PUT("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		var toUpdateBudget models.Budget
		if err := c.Bind(&toUpdateBudget); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		}

		props := database.BudgetProps{
			Year:   intPtr(toUpdateBudget.Year),
			Month:  intPtr(toUpdateBudget.Month),
			AreaId: intPtr(toUpdateBudget.AreaID),
			Amount: intPtr(toUpdateBudget.Amount),
		}

		updatedBudget, err := database.UpdateBudget(id, props)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedBudget)
	})
}

// DeleteBudget defines the route for deleting a budget by its ID.
func DeleteBudget(budgetGroup *echo.Group) {
	budgetGroup.DELETE("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		err = database.DeleteBudget(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Orçamento excluído com sucesso"})
	})
}
