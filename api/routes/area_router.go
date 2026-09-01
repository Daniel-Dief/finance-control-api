package routes

import (
	"finance-control-api/database"
	"finance-control-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// ListAreas defines the route for listing areas.
func ListAreas(areaGroup *echo.Group) {
	areaGroup.GET("/list", func(c *echo.Context) error {
		name := c.QueryParam("name")

		areas, err := database.ListAreas(name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if len(areas) == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Nenhuma área encontrada"})
		}

		return c.JSON(http.StatusOK, areas)
	})
}

// GetAreaByID defines the route for retrieving an area by its ID.
func GetAreaByID(areaGroup *echo.Group) {
	areaGroup.GET("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		area, err := database.GetAreaByID(id)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if area.ID == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Área não encontrada"})
		}

		return c.JSON(http.StatusOK, area)
	})
}

// CreateArea defines the route for creating a new area.
func CreateArea(areaGroup *echo.Group) {
	areaGroup.POST("/create", func(c *echo.Context) error {
		var toCreateArea models.Area

		if err := c.Bind(&toCreateArea); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		}

		if toCreateArea.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O nome da área é obrigatório"})
		}

		createdArea, err := database.CreateArea(toCreateArea.Name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, createdArea)
	})
}

// UpdateArea defines the route for updating an existing area.
func UpdateArea(areaGroup *echo.Group) {
	areaGroup.PUT("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		var toUpdateArea models.Area
		if err := c.Bind(&toUpdateArea); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		}

		if toUpdateArea.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O nome da área é obrigatório"})
		}

		updatedArea, err := database.UpdateArea(id, toUpdateArea.Name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedArea)
	})
}

// DeleteArea defines the route for deleting an area by its ID.
func DeleteArea(areaGroup *echo.Group) {
	areaGroup.DELETE("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		err = database.DeleteArea(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Área excluída com sucesso"})
	})
}
