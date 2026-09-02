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

func TestListTransactions(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/transactions/list", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusNoContent)
}

func TestListTransactionsWithFilters(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Transaction Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Transaction Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	_, err = createTestTransaction("2024-06-15", 100, category.ID, area.ID, "income")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/transactions/list?type=income", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var transactions []models.Transaction
	err = json.Unmarshal(rec.Body.Bytes(), &transactions)
	assert.NoError(t, err)
	assert.NotEmpty(t, transactions)

	for _, tr := range transactions {
		database.DeleteTransaction(tr.ID)
	}
	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestGetTransactionByID(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Transaction Get Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Transaction Get Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	transaction, err := createTestTransaction("2024-06-15", 150, category.ID, area.ID, "expense")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transactions/%d", transaction.ID), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Transaction
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID, response.ID)

	database.DeleteTransaction(transaction.ID)
	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestGetTransactionByIDNotFound(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/transactions/99999", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetTransactionByIDInvalid(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/transactions/abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTransaction(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Transaction Create Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Transaction Create Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	body := fmt.Sprintf(`{"date":"2024-06-15","amount":200,"category_id":%d,"area_id":%d,"type":"income"}`, category.ID, area.ID)
	req := httptest.NewRequest(http.MethodPost, "/transactions/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["id"])
	assert.Equal(t, "2024-06-15T00:00:00Z", response["date"])

	if id, ok := response["id"].(float64); ok {
		database.DeleteTransaction(int(id))
	}
	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestCreateTransactionInvalidBody(t *testing.T) {
	e := setupEcho()

	body := `{"invalid": "json"}`
	req := httptest.NewRequest(http.MethodPost, "/transactions/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTransactionMissingFields(t *testing.T) {
	e := setupEcho()

	body := `{"date":"2024-06-15"}`
	req := httptest.NewRequest(http.MethodPost, "/transactions/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateTransactionInvalidType(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Transaction Invalid Type Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Transaction Invalid Type Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	body := fmt.Sprintf(`{"date":"2024-06-15","amount":100,"category_id":%d,"area_id":%d,"type":"invalid"}`, category.ID, area.ID)
	req := httptest.NewRequest(http.MethodPost, "/transactions/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestCreateTransactionInvalidDate(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Transaction Invalid Date Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Transaction Invalid Date Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	body := fmt.Sprintf(`{"date":"invalid-date","amount":100,"category_id":%d,"area_id":%d,"type":"income"}`, category.ID, area.ID)
	req := httptest.NewRequest(http.MethodPost, "/transactions/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestUpdateTransaction(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Transaction Update Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Transaction Update Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	transaction, err := createTestTransaction("2024-06-15", 100, category.ID, area.ID, "income")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	body := `{"amount":250}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/transactions/%d", transaction.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Transaction
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 250, response.Amount)

	database.DeleteTransaction(transaction.ID)
	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestUpdateTransactionNotFound(t *testing.T) {
	e := setupEcho()

	body := `{"amount":100}`
	req := httptest.NewRequest(http.MethodPut, "/transactions/99999", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateTransactionInvalidID(t *testing.T) {
	e := setupEcho()

	body := `{"amount":100}`
	req := httptest.NewRequest(http.MethodPut, "/transactions/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteTransaction(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Transaction Delete Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Transaction Delete Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	transaction, err := createTestTransaction("2024-06-15", 100, category.ID, area.ID, "expense")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/transactions/%d", transaction.ID), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Transação excluída com sucesso", response["info"])

	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestDeleteTransactionNotFound(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/transactions/99999", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteTransactionInvalidID(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/transactions/abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}