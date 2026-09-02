package routes

import (
	"errors"
	"finance-control-api/database"
	"net/http"

	"github.com/labstack/echo/v5"
)

// TransactionQueryParams defines the query parameters accepted by ListTransactions.
// Pointer fields stay nil when the parameter is omitted.
type TransactionQueryParams struct {
	Type       *string `query:"type" validate:"omitempty,oneof=income expense"`
	CategoryID *int    `query:"category_id" validate:"omitempty,gt=0"`
	AreaID     *int    `query:"area_id" validate:"omitempty,gt=0"`
	From       *string `query:"from" validate:"omitempty,datetime=2006-01-02"`
	To         *string `query:"to" validate:"omitempty,datetime=2006-01-02"`
}

// TransactionBody defines the JSON body accepted by CreateTransaction.
type TransactionBody struct {
	Date       string `json:"date" validate:"required,datetime=2006-01-02"`
	Amount     int    `json:"amount" validate:"required"`
	CategoryID int    `json:"category_id" validate:"required,gt=0"`
	AreaID     int    `json:"area_id" validate:"required,gt=0"`
	Type       string `json:"type" validate:"required,oneof=income expense"`
}

// TransactionUpdateBody defines the JSON body accepted by UpdateTransaction. All
// fields are optional: the update only applies the fields that were provided.
type TransactionUpdateBody struct {
	Date       *string `json:"date" validate:"omitempty,datetime=2006-01-02"`
	Amount     *int    `json:"amount" validate:"omitempty"`
	CategoryID *int    `json:"category_id" validate:"omitempty,gt=0"`
	AreaID     *int    `json:"area_id" validate:"omitempty,gt=0"`
	Type       *string `json:"type" validate:"omitempty,oneof=income expense"`
}

// ListTransactions defines the route for listing transactions.
func ListTransactions(transactionGroup *echo.Group) {
	transactionGroup.GET("/list", func(c *echo.Context) error {
		var params TransactionQueryParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		filters := database.TransactionFilters{
			Type:       params.Type,
			CategoryID: params.CategoryID,
			AreaID:     params.AreaID,
			From:       params.From,
			To:         params.To,
		}

		transactions, err := database.ListTransactions(filters)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if len(transactions) == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Nenhuma transação encontrada"})
		}

		return c.JSON(http.StatusOK, transactions)
	})
}

// GetTransactionByID defines the route for retrieving a transaction by its ID.
func GetTransactionByID(transactionGroup *echo.Group) {
	transactionGroup.GET("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		transaction, err := database.GetTransactionByID(params.ID)

		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Transação não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, transaction)
	})
}

// CreateTransaction defines the route for creating a new transaction.
func CreateTransaction(transactionGroup *echo.Group) {
	transactionGroup.POST("/create", func(c *echo.Context) error {
		var body TransactionBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		props := database.TransactionProps{
			Date:       &body.Date,
			Amount:     &body.Amount,
			CategoryID: &body.CategoryID,
			AreaID:     &body.AreaID,
			Type:       &body.Type,
		}

		createdTransaction, err := database.CreateTransaction(props)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, createdTransaction)
	})
}

// UpdateTransaction defines the route for updating an existing transaction.
func UpdateTransaction(transactionGroup *echo.Group) {
	transactionGroup.PUT("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		var body TransactionUpdateBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		props := database.TransactionProps{
			Date:       body.Date,
			Amount:     body.Amount,
			CategoryID: body.CategoryID,
			AreaID:     body.AreaID,
			Type:       body.Type,
		}

		updatedTransaction, err := database.UpdateTransaction(params.ID, props)

		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Transação não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedTransaction)
	})
}

// DeleteTransaction defines the route for deleting a transaction by its ID.
func DeleteTransaction(transactionGroup *echo.Group) {
	transactionGroup.DELETE("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		err := database.DeleteTransaction(params.ID)
		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Transação não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Transação excluída com sucesso"})
	})
}
