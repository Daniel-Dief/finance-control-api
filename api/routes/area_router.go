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

// ListAreas defines the route for listing areas.
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

// GetAreaByID defines the route for retrieving an area by its ID.
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

// CreateArea defines the route for creating a new area.
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

// UpdateArea defines the route for updating an existing area.
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

// DeleteArea defines the route for deleting an area by its ID.
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
