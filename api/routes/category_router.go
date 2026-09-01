package routes

import (
	"finance-control-api/database"
	"finance-control-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// ListCategories defines the route for listing categories.
func ListCategories(categoryGroup *echo.Group) {
	categoryGroup.GET("/list", func(c *echo.Context) error {
		name := c.QueryParam("name")

		var namePtr *string
		if name != "" {
			namePtr = &name
		}

		categories, err := database.ListCategories(namePtr)

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
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		category, err := database.GetCategoryByID(id)

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
		var toCreateCategory models.Category

		if err := c.Bind(&toCreateCategory); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		}

		if toCreateCategory.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O nome da categoria é obrigatório"})
		}

		createdCategory, err := database.CreateCategory(toCreateCategory.Name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, createdCategory)
	})
}

// UpdateCategory defines the route for updating an existing category.
func UpdateCategory(categoryGroup *echo.Group) {
	categoryGroup.PUT("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		var toUpdateCategory models.Category
		if err := c.Bind(&toUpdateCategory); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Erro ao processar o corpo da requisição."})
		}

		if toUpdateCategory.Name == "" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "O nome da categoria é obrigatório"})
		}

		updatedCategory, err := database.UpdateCategory(id, toUpdateCategory.Name)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, updatedCategory)
	})
}

// DeleteCategory defines the route for deleting a category by its ID.
func DeleteCategory(categoryGroup *echo.Group) {
	categoryGroup.DELETE("/:id", func(c *echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "ID inválido"})
		}

		err = database.DeleteCategory(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{"info": "Categoria excluída com sucesso"})
	})
}
