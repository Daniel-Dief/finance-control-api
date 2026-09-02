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
	Page       *int    `query:"page" validate:"omitempty,min=1"`
	Limit      *int    `query:"limit" validate:"omitempty,min=1,max=100"`
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

// ListTransactions godoc
//
//	@Summary		Listar transações
//	@Description	Retorna a lista de transações, com filtros opcionais por tipo, categoria, área e período.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			type			query		string	false	"Filtro por tipo (income ou expense)"	Enums(income, expense)
//	@Param			category_id		query		int		false	"Filtro por ID da categoria"
//	@Param			area_id			query		int		false	"Filtro por ID da área"
//	@Param			from			query		string	false	"Data inicial (AAAA-MM-DD)"				Format(2006-01-02)
//	@Param			to				query		string	false	"Data final (AAAA-MM-DD)"					Format(2006-01-02)
//	@Param			page			query		int		false	"Número da página (padrão 1)"
//	@Param			limit			query		int		false	"Quantidade por página (padrão 50, máximo 100)"
//	@Success		200				{object}	database.PaginatedResult[models.Transaction]
//	@Failure		500				{object}	map[string]string
//	@Router			/transactions/list [get]
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

		pag := parsePagination(params.Page, params.Limit)

		transactions, err := database.ListTransactions(filters, pag)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if len(transactions.Data) == 0 && transactions.Total == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Nenhuma transação encontrada"})
		}

		return c.JSON(http.StatusOK, transactions)
	})
}

// GetTransactionByID godoc
//
//	@Summary		Obter transação por ID
//	@Description	Retorna os dados de uma transação a partir do seu identificador.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID da transação"
//	@Success		200	{object}	models.Transaction
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/transactions/{id} [get]
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

// CreateTransaction godoc
//
//	@Summary		Criar transação
//	@Description	Cria uma nova transação (receita ou despesa) a partir dos dados informados.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		TransactionBody	true	"Dados da nova transação"
//	@Success		201		{object}	models.Transaction
//	@Failure		400		{object}		map[string]interface{}
//	@Failure		500		{object}		map[string]string
//	@Router			/transactions/create [post]
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

// UpdateTransaction godoc
//
//	@Summary		Atualizar transação
//	@Description	Atualiza os campos de uma transação existente. Apenas os campos informados serão alterados.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"ID da transação"
//	@Param			body	body		TransactionUpdateBody	true	"Campos a atualizar"
//	@Success		200		{object}	models.Transaction
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/transactions/{id} [put]
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

// DeleteTransaction godoc
//
//	@Summary		Excluir transação
//	@Description	Remove uma transação a partir do seu identificador.
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID da transação"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/transactions/{id} [delete]
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
