package routes

import (
	"finance-control-api/database"
	"finance-control-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// ListTransactions defines the route for listing transactions.
func ListTransactions(transactionGroup *echo.Group) {
	transactionGroup.GET("/list", func(c *echo.Context) error {
		filters := database.TransactionFilters{
			Type:       stringPtr(c.QueryParam("type")),
			CategoryID: parseIntPtr(c.QueryParam("category_id")),
			AreaID:     parseIntPtr(c.QueryParam("area_id")),
			From:       stringPtr(c.QueryParam("from")),
			To:         stringPtr(c.QueryParam("to")),
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
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		transaction, err := database.GetTransactionByID(id)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if transaction.ID == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Transação não encontrada"})
		}

		return c.JSON(http.StatusOK, transaction)
	})
}

// CreateTransaction defines the route for creating a new transaction.
func CreateTransaction(transactionGroup *echo.Group) {
	transactionGroup.POST("/create", func(c *echo.Context) error {
		var toCreateTransaction models.Transaction

		if err := c.Bind(&toCreateTransaction); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		}

		if toCreateTransaction.Type == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O tipo da transação é obrigatório"})
		}
		if toCreateTransaction.CategoryID == 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O id da categoria é obrigatório"})
		}
		if toCreateTransaction.AreaID == 0 {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O id da área é obrigatório"})
		}

		props := database.TransactionProps{
			Date:       &toCreateTransaction.Date,
			Amount:     &toCreateTransaction.Amount,
			CategoryID: &toCreateTransaction.CategoryID,
			AreaID:     &toCreateTransaction.AreaID,
			Type:       &toCreateTransaction.Type,
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
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		var toUpdateTransaction models.Transaction
		if err := c.Bind(&toUpdateTransaction); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		}

		props := database.TransactionProps{
			Date:       stringPtr(toUpdateTransaction.Date),
			Amount:     intPtr(toUpdateTransaction.Amount),
			CategoryID: intPtr(toUpdateTransaction.CategoryID),
			AreaID:     intPtr(toUpdateTransaction.AreaID),
			Type:       stringPtr(toUpdateTransaction.Type),
		}

		updatedTransaction, err := database.UpdateTransaction(id, props)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedTransaction)
	})
}

// DeleteTransaction defines the route for deleting a transaction by its ID.
func DeleteTransaction(transactionGroup *echo.Group) {
	transactionGroup.DELETE("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		err = database.DeleteTransaction(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Transação excluída com sucesso"})
	})
}
