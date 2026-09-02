package routes

import (
	"errors"
	"finance-control-api/database"
	"net/http"

	"github.com/labstack/echo/v5"
)

// AreaQueryParams defines the query parameters accepted by ListAreas.
type AreaQueryParams struct {
	Name string `query:"name" validate:"omitempty"`
}

// AreaBody defines the JSON body accepted by CreateArea and UpdateArea.
type AreaBody struct {
	Name string `json:"name" validate:"required"`
}

// ListAreas godoc
//
//	@Summary		Listar áreas
//	@Description	Retorna a lista de todas as áreas cadastradas, com filtro opcional por nome.
//	@Tags			areas
//	@Accept			json
//	@Produce		json
//	@Param			name	query		string	false	"Filtro por nome da área"
//	@Success		200		{array}		models.Area
//	@Success		204		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/areas/list [get]
func ListAreas(areaGroup *echo.Group) {
	areaGroup.GET("/list", func(c *echo.Context) error {
		var params AreaQueryParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		areas, err := database.ListAreas(params.Name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if len(areas) == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Nenhuma área encontrada"})
		}

		return c.JSON(http.StatusOK, areas)
	})
}

// GetAreaByID godoc
//
//	@Summary		Obter área por ID
//	@Description	Retorna os dados de uma área a partir do seu identificador.
//	@Tags			areas
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID da área"
//	@Success		200	{object}	models.Area
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/areas/{id} [get]
func GetAreaByID(areaGroup *echo.Group) {
	areaGroup.GET("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		area, err := database.GetAreaByID(params.ID)

		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Área não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, area)
	})
}

// CreateArea godoc
//
//	@Summary		Criar área
//	@Description	Cria uma nova área a partir do nome informado no corpo da requisição.
//	@Tags			areas
//	@Accept			json
//	@Produce		json
//	@Param			body	body		AreaBody	true	"Dados da nova área"
//	@Success		201		{object}	models.Area
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]string
//	@Router			/areas/create [post]
func CreateArea(areaGroup *echo.Group) {
	areaGroup.POST("/create", func(c *echo.Context) error {
		var body AreaBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		createdArea, err := database.CreateArea(body.Name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, createdArea)
	})
}

// UpdateArea godoc
//
//	@Summary		Atualizar área
//	@Description	Atualiza o nome de uma área existente identificada pelo ID.
//	@Tags			areas
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"ID da área"
//	@Param			body	body		AreaBody	true	"Novos dados da área"
//	@Success		200		{object}	models.Area
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/areas/{id} [put]
func UpdateArea(areaGroup *echo.Group) {
	areaGroup.PUT("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		var body AreaBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		updatedArea, err := database.UpdateArea(params.ID, body.Name)

		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Área não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedArea)
	})
}

// DeleteArea godoc
//
//	@Summary		Excluir área
//	@Description	Remove uma área a partir do seu identificador.
//	@Tags			areas
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"ID da área"
//	@Success		200	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/areas/{id} [delete]
func DeleteArea(areaGroup *echo.Group) {
	areaGroup.DELETE("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		err := database.DeleteArea(params.ID)
		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Área não encontrada"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Área excluída com sucesso"})
	})
}
