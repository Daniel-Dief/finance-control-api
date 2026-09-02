package routes

import (
	"errors"
	"finance-control-api/database"
	"net/http"

	"github.com/labstack/echo/v5"
)

// BudgetQueryParams defines the query parameters accepted by ListBudgets.
// Pointer fields stay nil when the parameter is omitted.
type BudgetQueryParams struct {
	Year   *int `query:"year" validate:"omitempty,min=1"`
	Month  *int `query:"month" validate:"omitempty,min=1,max=12"`
	AreaID *int `query:"area_id" validate:"omitempty,gt=0"`
}

// BudgetBody defines the JSON body accepted by CreateBudget.
type BudgetBody struct {
	Year   int `json:"year" validate:"required,min=1"`
	Month  int `json:"month" validate:"required,min=1,max=12"`
	AreaID int `json:"area_id" validate:"required,gt=0"`
	Amount int `json:"amount"`
}

// BudgetUpdateBody defines the JSON body accepted by UpdateBudget. All fields
// are optional: the update only applies the fields that were provided.
type BudgetUpdateBody struct {
	Year   *int `json:"year" validate:"omitempty,min=1"`
	Month  *int `json:"month" validate:"omitempty,min=1,max=12"`
	AreaID *int `json:"area_id" validate:"omitempty,gt=0"`
	Amount *int `json:"amount" validate:"omitempty"`
}

// ListBudgets defines the route for listing budgets.
func ListBudgets(budgetGroup *echo.Group) {
	budgetGroup.GET("/list", func(c *echo.Context) error {
		var params BudgetQueryParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		props := database.BudgetProps{
			Year:   params.Year,
			Month:  params.Month,
			AreaId: params.AreaID,
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
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		budget, err := database.GetBudgetByID(params.ID)

		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Orçamento não encontrado"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, budget)
	})
}

// CreateBudget defines the route for creating a new budget.
func CreateBudget(budgetGroup *echo.Group) {
	budgetGroup.POST("/create", func(c *echo.Context) error {
		var body BudgetBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		props := database.BudgetProps{
			Year:   &body.Year,
			Month:  &body.Month,
			AreaId: &body.AreaID,
			Amount: &body.Amount,
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
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		var body BudgetUpdateBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		props := database.BudgetProps{
			Year:   body.Year,
			Month:  body.Month,
			AreaId: body.AreaID,
			Amount: body.Amount,
		}

		updatedBudget, err := database.UpdateBudget(params.ID, props)

		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Orçamento não encontrado"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedBudget)
	})
}

// DeleteBudget defines the route for deleting a budget by its ID.
func DeleteBudget(budgetGroup *echo.Group) {
	budgetGroup.DELETE("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		err := database.DeleteBudget(params.ID)
		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Orçamento não encontrado"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Orçamento excluído com sucesso"})
	})
}
