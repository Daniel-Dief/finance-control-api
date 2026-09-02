package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"finance-control-api/database"
	"finance-control-api/models"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestListCategories(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/categories/list", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusNoContent)
}

func TestListCategoriesWithNameFilter(t *testing.T) {
	e := setupEcho()

	_, err := createTestCategory("Test Category for List")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/categories/list?name=Test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var categories []models.Category
	err = json.Unmarshal(rec.Body.Bytes(), &categories)
	assert.NoError(t, err)
	assert.NotEmpty(t, categories)

	for _, c := range categories {
		database.DeleteCategory(c.ID)
	}
}

func TestGetCategoryByID(t *testing.T) {
	e := setupEcho()

	category, err := createTestCategory("Category for Get Test")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/categories/%d", category.ID), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Category
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, category.ID, response.ID)

	database.DeleteCategory(category.ID)
}

func TestGetCategoryByIDNotFound(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/categories/99999", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetCategoryByIDInvalid(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/categories/abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCategory(t *testing.T) {
	e := setupEcho()

	body := `{"name":"New Test Category"}`
	req := httptest.NewRequest(http.MethodPost, "/categories/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["id"])
	assert.Equal(t, "New Test Category", response["name"])

	if id, ok := response["id"].(float64); ok {
		database.DeleteCategory(int(id))
	}
}

func TestCreateCategoryInvalidBody(t *testing.T) {
	e := setupEcho()

	body := `{"invalid": "json"}`
	req := httptest.NewRequest(http.MethodPost, "/categories/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateCategoryEmptyName(t *testing.T) {
	e := setupEcho()

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPost, "/categories/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateCategory(t *testing.T) {
	e := setupEcho()

	category, err := createTestCategory("Category to Update")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	body := `{"name":"Updated Category Name"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/categories/%d", category.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Category
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Category Name", response.Name)

	database.DeleteCategory(category.ID)
}

func TestUpdateCategoryNotFound(t *testing.T) {
	e := setupEcho()

	body := `{"name":"Updated Name"}`
	req := httptest.NewRequest(http.MethodPut, "/categories/99999", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateCategoryInvalidID(t *testing.T) {
	e := setupEcho()

	body := `{"name":"Updated Name"}`
	req := httptest.NewRequest(http.MethodPut, "/categories/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteCategory(t *testing.T) {
	e := setupEcho()

	category, err := createTestCategory("Category to Delete")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/categories/%d", category.ID), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Categoria excluída com sucesso", response["info"])

	database.DeleteCategory(category.ID)
}

func TestDeleteCategoryNotFound(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/categories/99999", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteCategoryInvalidID(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/categories/abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}