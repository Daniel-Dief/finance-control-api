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

// TestGetAreaByIDAfterDelete ensures a soft-deleted area is treated as non-existent.
func TestGetAreaByIDAfterDelete(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Area Get")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	err = database.DeleteArea(testUserID,area.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/areas/%d", area.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestUpdateAreaAfterDelete ensures a soft-deleted area cannot be updated.
func TestUpdateAreaAfterDelete(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Area Update")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	err = database.DeleteArea(testUserID,area.ID)
	assert.NoError(t, err)

	body := `{"name":"Should Not Update"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/areas/%d", area.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestDeleteAreaTwice ensures deleting an already soft-deleted area returns 404.
func TestDeleteAreaTwice(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Area Twice")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/areas/%d", area.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/areas/%d", area.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestListAreasExcludesDeleted ensures a soft-deleted area no longer shows in the list.
func TestListAreasExcludesDeleted(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Area List")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	err = database.DeleteArea(testUserID,area.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/areas/list?name=Soft%20Delete%20Area%20List", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// TestGetCategoryByIDAfterDelete ensures a soft-deleted category is treated as non-existent.
func TestGetCategoryByIDAfterDelete(t *testing.T) {
	e := setupEcho()

	category, err := createTestCategory("Soft Delete Category Get")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	err = database.DeleteCategory(testUserID,category.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/categories/%d", category.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestUpdateCategoryAfterDelete ensures a soft-deleted category cannot be updated.
func TestUpdateCategoryAfterDelete(t *testing.T) {
	e := setupEcho()

	category, err := createTestCategory("Soft Delete Category Update")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	err = database.DeleteCategory(testUserID,category.ID)
	assert.NoError(t, err)

	body := `{"name":"Should Not Update"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/categories/%d", category.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestDeleteCategoryTwice ensures deleting an already soft-deleted category returns 404.
func TestDeleteCategoryTwice(t *testing.T) {
	e := setupEcho()

	category, err := createTestCategory("Soft Delete Category Twice")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/categories/%d", category.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/categories/%d", category.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestListCategoriesExcludesDeleted ensures a soft-deleted category no longer shows in the list.
func TestListCategoriesExcludesDeleted(t *testing.T) {
	e := setupEcho()

	category, err := createTestCategory("Soft Delete Category List")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	err = database.DeleteCategory(testUserID,category.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/categories/list?name=Soft%20Delete%20Category%20List", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// TestGetBudgetByIDAfterDelete ensures a soft-deleted budget is treated as non-existent.
func TestGetBudgetByIDAfterDelete(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Budget Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,area.ID)

	budget, err := createTestBudget(2024, 2, area.ID, 1000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	err = database.DeleteBudget(testUserID,budget.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/budgets/%d", budget.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestUpdateBudgetAfterDelete ensures a soft-deleted budget cannot be updated.
func TestUpdateBudgetAfterDelete(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Budget Update Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,area.ID)

	budget, err := createTestBudget(2024, 3, area.ID, 1000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	err = database.DeleteBudget(testUserID,budget.ID)
	assert.NoError(t, err)

	body := `{"amount":2000}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/budgets/%d", budget.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestDeleteBudgetTwice ensures deleting an already soft-deleted budget returns 404.
func TestDeleteBudgetTwice(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Budget Twice Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,area.ID)

	budget, err := createTestBudget(2024, 4, area.ID, 1000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/budgets/%d", budget.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/budgets/%d", budget.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestListBudgetsExcludesDeleted ensures a soft-deleted budget no longer shows in the list.
func TestListBudgetsExcludesDeleted(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Budget List Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,area.ID)

	budget, err := createTestBudget(2024, 5, area.ID, 1000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	err = database.DeleteBudget(testUserID,budget.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/budgets/list?year=2024&month=5&area_id=%d", area.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// TestGetTransactionByIDAfterDelete ensures a soft-deleted transaction is treated as non-existent.
func TestGetTransactionByIDAfterDelete(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Transaction Get Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,area.ID)

	category, err := createTestCategory("Soft Delete Transaction Get Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	defer database.DeleteCategory(testUserID,category.ID)

	transaction, err := createTestTransaction("2024-06-15", 100, category.ID, area.ID, "income")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	err = database.DeleteTransaction(testUserID,transaction.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%d", transaction.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestUpdateTransactionAfterDelete ensures a soft-deleted transaction cannot be updated.
func TestUpdateTransactionAfterDelete(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Transaction Update Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,area.ID)

	category, err := createTestCategory("Soft Delete Transaction Update Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	defer database.DeleteCategory(testUserID,category.ID)

	transaction, err := createTestTransaction("2024-06-15", 100, category.ID, area.ID, "income")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	err = database.DeleteTransaction(testUserID,transaction.ID)
	assert.NoError(t, err)

	body := `{"amount":250}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/transactions/%d", transaction.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestDeleteTransactionTwice ensures deleting an already soft-deleted transaction returns 404.
func TestDeleteTransactionTwice(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Transaction Twice Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,area.ID)

	category, err := createTestCategory("Soft Delete Transaction Twice Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	defer database.DeleteCategory(testUserID,category.ID)

	transaction, err := createTestTransaction("2024-06-15", 100, category.ID, area.ID, "expense")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/transactions/%d", transaction.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/transactions/%d", transaction.ID), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// TestListTransactionsExcludesDeleted ensures a soft-deleted transaction no longer shows in the list.
func TestListTransactionsExcludesDeleted(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Soft Delete Transaction List Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,area.ID)

	category, err := createTestCategory("Soft Delete Transaction List Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}
	defer database.DeleteCategory(testUserID,category.ID)

	transaction, err := createTestTransaction("2024-06-15", 100, category.ID, area.ID, "income")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	err = database.DeleteTransaction(testUserID,transaction.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/list?type=income&category_id=%d", category.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// TestCreateAndVerifyListContainsActive verifies active records are still listed after a soft-delete elsewhere.
func TestCreateAndVerifyListContainsActive(t *testing.T) {
	e := setupEcho()

	active, err := createTestArea("Active Area Survives")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,active.ID)

	deleted, err := createTestArea("Deleted Area Hidden")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	defer database.DeleteArea(testUserID,deleted.ID)

	err = database.DeleteArea(testUserID,deleted.ID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/areas/list?name=Survives", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var areas database.PaginatedResult[models.Area]
	err = json.Unmarshal(rec.Body.Bytes(), &areas)
	assert.NoError(t, err)
	assert.Len(t, areas.Data, 1)
	assert.Equal(t, active.ID, areas.Data[0].ID)
}
