package routes

import (
	"finance-control-api/database"
	"net/http"

	"github.com/labstack/echo/v5"
)

// CategoryQueryParams defines the query parameters accepted by ListCategories.
// Pointer fields stay nil when the parameter is omitted.
type CategoryQueryParams struct {
	Name *string `query:"name" validate:"omitempty"`
}

// CategoryBody defines the JSON body accepted by CreateCategory and UpdateCategory.
type CategoryBody struct {
	Name string `json:"name" validate:"required"`
}

// ListCategories defines the route for listing categories.
func ListCategories(categoryGroup *echo.Group) {
	categoryGroup.GET("/list", func(c *echo.Context) error {
		var params CategoryQueryParams
		if !bindAndValidateQuery(c, &params) {
			return nil
		}

		categories, err := database.ListCategories(params.Name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if len(categories) == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Nenhuma categoria encontrada"})
		}

		return c.JSON(http.StatusOK, categories)
	})
}

// GetCategoryByID defines the route for retrieving a category by its ID.
func GetCategoryByID(categoryGroup *echo.Group) {
	categoryGroup.GET("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		category, err := database.GetCategoryByID(params.ID)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}
		if category.ID == 0 {
			return c.JSON(http.StatusNoContent, map[string]interface{}{"info": "Categoria não encontrada"})
		}

		return c.JSON(http.StatusOK, category)
	})
}

// CreateCategory defines the route for creating a new category.
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

// UpdateCategory defines the route for updating an existing category.
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

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedCategory)
	})
}

// DeleteCategory defines the route for deleting a category by its ID.
func DeleteCategory(categoryGroup *echo.Group) {
	categoryGroup.DELETE("/:id", func(c *echo.Context) error {
		var params PathIDParams
		if !bindAndValidatePath(c, &params) {
			return nil
		}

		err := database.DeleteCategory(params.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Categoria excluída com sucesso"})
	})
}
