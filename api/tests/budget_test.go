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

func TestListBudgets(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/budgets/list", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusNoContent)
}

func TestListBudgetsWithFilters(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Budget Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	_, err = createTestBudget(2024, 1, area.ID, 1000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/budgets/list?year=2024&month=1", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var budgets database.PaginatedResult[models.Budget]
	err = json.Unmarshal(rec.Body.Bytes(), &budgets)
	assert.NoError(t, err)
	assert.NotEmpty(t, budgets.Data)

	for _, b := range budgets.Data {
		database.DeleteBudget(b.ID)
	}
	database.DeleteArea(area.ID)
}

func TestGetBudgetByID(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Budget Get Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	budget, err := createTestBudget(2024, 6, area.ID, 5000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/budgets/%d", budget.ID), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Budget
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, budget.ID, response.ID)

	database.DeleteBudget(budget.ID)
	database.DeleteArea(area.ID)
}

func TestGetBudgetByIDNotFound(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/budgets/99999", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetBudgetByIDInvalid(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/budgets/abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateBudget(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Budget Create Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	body := fmt.Sprintf(`{"year":2024,"month":12,"area_id":%d,"amount":10000}`, area.ID)
	req := httptest.NewRequest(http.MethodPost, "/budgets/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["id"])
	assert.Equal(t, float64(2024), response["year"])
	assert.Equal(t, float64(12), response["month"])

	if id, ok := response["id"].(float64); ok {
		database.DeleteBudget(int(id))
	}
	database.DeleteArea(area.ID)
}

func TestCreateBudgetInvalidBody(t *testing.T) {
	e := setupEcho()

	body := `{"invalid": "json"}`
	req := httptest.NewRequest(http.MethodPost, "/budgets/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateBudgetMissingFields(t *testing.T) {
	e := setupEcho()

	body := `{"year":2024}`
	req := httptest.NewRequest(http.MethodPost, "/budgets/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateBudgetInvalidMonth(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Budget Invalid Month Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	body := fmt.Sprintf(`{"year":2024,"month":13,"area_id":%d,"amount":1000}`, area.ID)
	req := httptest.NewRequest(http.MethodPost, "/budgets/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	database.DeleteArea(area.ID)
}

func TestUpdateBudget(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Budget Update Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	budget, err := createTestBudget(2024, 6, area.ID, 5000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	body := `{"amount":7500}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/budgets/%d", budget.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Budget
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 7500, response.Amount)

	database.DeleteBudget(budget.ID)
	database.DeleteArea(area.ID)
}

func TestUpdateBudgetNotFound(t *testing.T) {
	e := setupEcho()

	body := `{"amount":1000}`
	req := httptest.NewRequest(http.MethodPut, "/budgets/99999", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateBudgetInvalidID(t *testing.T) {
	e := setupEcho()

	body := `{"amount":1000}`
	req := httptest.NewRequest(http.MethodPut, "/budgets/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteBudget(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Budget Delete Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	budget, err := createTestBudget(2024, 6, area.ID, 5000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/budgets/%d", budget.ID), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Orçamento excluído com sucesso", response["info"])

	database.DeleteBudget(budget.ID)
	database.DeleteArea(area.ID)
}

func TestDeleteBudgetNotFound(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/budgets/99999", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteBudgetInvalidID(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/budgets/abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}