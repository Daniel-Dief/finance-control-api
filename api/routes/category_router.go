package routes

import (
	"errors"
	"finance-control-api/database"
	"net/http"

	"github.com/labstack/echo/v5"
)

// CategoryQueryParams defines the query parameters accepted by ListCategories.
// Pointer fields stay nil when the parameter is omitted.
type CategoryQueryParams struct {
	Name  *string `query:"name" validate:"omitempty"`
	Page  *int    `query:"page" validate:"omitempty,min=1"`
	Limit *int    `query:"limit" validate:"omitempty,min=1,max=100"`
}

// CategoryBody defines the JSON body accepted by CreateCategory and UpdateCategory.
type CategoryBody struct {
	Name string `json:"name" validate:"required"`
}

// ListCategories godoc
//
//	@Summary		Listar categorias
//	@Description	Retorna a lista de categorias cadastradas, com filtro opcional por nome.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"Filtro por nome da categoria"
//	@Param			page	query		int		false	"Número da página (padrão 1)"
//	@Param			limit	query		int		false	"Quantidade por página (padrão 50, máximo 100)"
//	@Success		200		{object}	database.PaginatedResult[models.Category]
//	@Failure		500		{object}	map[string]string
//	@Router			/categories/list [get]
func ListCategories(categoryGroup *echo.Group) {
	categoryGroup.GET("/list", func(c *echo.Context) error {
		var params CategoryQueryParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		pag := parsePagination(params.Page, params.Limit)

		categories, err := database.ListCategories(params.Name, pag)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if len(categories.Data) == 0 && categories.Total == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Nenhuma categoria encontrada"})
		}

		return c.JSON(http.StatusOK, categories)
	})
}

// GetCategoryByID godoc
//
//	@Summary		Obter categoria por ID
//	@Description	Retorna os dados de uma categoria a partir do seu identificador.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID da categoria"
//	@Success		200	{object}	models.Category
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/categories/{id} [get]
func GetCategoryByID(categoryGroup *echo.Group) {
	categoryGroup.GET("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		category, err := database.GetCategoryByID(params.ID)

		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Categoria não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, category)
	})
}

// CreateCategory godoc
//
//	@Summary		Criar categoria
//	@Description	Cria uma nova categoria a partir do nome informado no corpo da requisição.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CategoryBody	true	"Dados da nova categoria"
//	@Success		201		{object}	models.Category
//	@Failure		400		{object}		map[string]interface{}
//	@Failure		500		{object}		map[string]string
//	@Router			/categories/create [post]
func CreateCategory(categoryGroup *echo.Group) {
	categoryGroup.POST("/create", func(c *echo.Context) error {
		var body CategoryBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		createdCategory, err := database.CreateCategory(body.Name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, createdCategory)
	})
}

// UpdateCategory godoc
//
//	@Summary		Atualizar categoria
//	@Description	Atualiza o nome de uma categoria existente identificada pelo ID.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"ID da categoria"
//	@Param			body	body		CategoryBody	true	"Novos dados da categoria"
//	@Success		200		{object}	models.Category
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/categories/{id} [put]
func UpdateCategory(categoryGroup *echo.Group) {
	categoryGroup.PUT("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		var body CategoryBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		updatedCategory, err := database.UpdateCategory(params.ID, body.Name)

		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Categoria não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedCategory)
	})
}

// DeleteCategory godoc
//
//	@Summary		Excluir categoria
//	@Description	Remove uma categoria a partir do seu identificador.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID da categoria"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/categories/{id} [delete]
func DeleteCategory(categoryGroup *echo.Group) {
	categoryGroup.DELETE("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		err := database.DeleteCategory(params.ID)
		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Categoria não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Categoria excluída com sucesso"})
	})
}
