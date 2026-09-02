package tests

import (
	"finance-control-api/database"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestInvalidJSON(t *testing.T) {
	e := setupEcho()

	invalidJSON := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/areas/create", strings.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestMissingContentType(t *testing.T) {
	e := setupEcho()

	body := `{"name":"Test"}`
	req := httptest.NewRequest(http.MethodPost, "/areas/create", strings.NewReader(body))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestEmptyBody(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodPost, "/areas/create", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestInvalidIDFormat(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/areas/not-a-number", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestNegativeID(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/areas/-1", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestZeroID(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/areas/0", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestBudgetInvalidMonth(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Validation Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	body := `{"year":2024,"month":13,"area_id":1,"amount":1000}`
	req := httptest.NewRequest(http.MethodPost, "/budgets/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	database.DeleteArea(area.ID)
}

func TestBudgetInvalidYear(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Validation Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	body := `{"year":0,"month":1,"area_id":1,"amount":1000}`
	req := httptest.NewRequest(http.MethodPost, "/budgets/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	database.DeleteArea(area.ID)
}

func TestTransactionInvalidType(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Validation Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Validation Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	body := `{"date":"2024-06-15","amount":100,"category_id":1,"area_id":1,"type":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/transactions/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestTransactionInvalidDateFormat(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Validation Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	category, err := createTestCategory("Validation Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	body := `{"date":"15-06-2024","amount":100,"category_id":1,"area_id":1,"type":"income"}`
	req := httptest.NewRequest(http.MethodPost, "/transactions/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	database.DeleteCategory(category.ID)
	database.DeleteArea(area.ID)
}

func TestQueryParameterValidation(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/budgets/list?year=abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestQueryParameterInvalidMonth(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/budgets/list?month=13", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestQueryParameterInvalidAreaID(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/budgets/list?area_id=-1", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionQueryParameterInvalidType(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/transactions/list?type=invalid", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionQueryParameterInvalidDate(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/transactions/list?from=invalid-date", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTransactionQueryParameterInvalidCategoryID(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/transactions/list?category_id=abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}