package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"finance-control-api/models"

	"github.com/stretchr/testify/assert"
)

func TestGetKPIs(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("KPI Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	category, err := createTestCategory("KPI Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	_, err = createTestTransaction(today, 100000, category.ID, area.ID, "income")
	if err != nil {
		t.Fatalf("Failed to create test income: %v", err)
	}
	_, err = createTestTransaction(today, 60000, category.ID, area.ID, "expense")
	if err != nil {
		t.Fatalf("Failed to create test expense: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/reports/kpis", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var kpis models.ReportKPIs
	err = json.Unmarshal(rec.Body.Bytes(), &kpis)
	assert.NoError(t, err)
	assert.Equal(t, 100000, kpis.Income)
	assert.Equal(t, 60000, kpis.Expense)
	assert.Equal(t, 40000, kpis.Balance)
}

func TestGetLineChart(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Line Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	category, err := createTestCategory("Line Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	_, err = createTestTransaction(today, 15000, category.ID, area.ID, "expense")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/reports/line-chart?type=expense", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var chart models.LineChart
	err = json.Unmarshal(rec.Body.Bytes(), &chart)
	assert.NoError(t, err)
	assert.Equal(t, "expense", chart.Type)
	assert.Equal(t, 6, chart.HistoryMonths)
	assert.True(t, len(chart.Points) >= 28 && len(chart.Points) <= 31)
}

func TestGetRadarChart(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Radar Test Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	category, err := createTestCategory("Radar Test Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	today := time.Now().Format("2006-01-02")
	_, err = createTestTransaction(today, 20000, category.ID, area.ID, "expense")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/reports/radar-chart?type=expense&area_ids=%d", area.ID), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var chart models.RadarChart
	err = json.Unmarshal(rec.Body.Bytes(), &chart)
	assert.NoError(t, err)
	assert.NotEmpty(t, chart.Areas)
}

func TestGetBudgetReport(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Budget Report Area")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}
	category, err := createTestCategory("Budget Report Category")
	if err != nil {
		t.Fatalf("Failed to create test category: %v", err)
	}

	now := time.Now()
	budget, err := createTestBudget(now.Year(), int(now.Month()), area.ID, 100000)
	if err != nil {
		t.Fatalf("Failed to create test budget: %v", err)
	}

	today := now.Format("2006-01-02")
	_, err = createTestTransaction(today, 75000, category.ID, area.ID, "expense")
	if err != nil {
		t.Fatalf("Failed to create test transaction: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/reports/budgets", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var report models.BudgetReport
	err = json.Unmarshal(rec.Body.Bytes(), &report)
	assert.NoError(t, err)
	assert.NotEmpty(t, report.Budgets)

	for _, b := range report.Budgets {
		if b.ID == budget.ID {
			assert.Equal(t, 75000, b.Spent)
			assert.Equal(t, 25000, b.Remaining)
			assert.Equal(t, 75.0, b.Percentage)
			assert.False(t, b.Exceeded)
			break
		}
	}
}

func TestReportsPeriodValidation(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/reports/kpis?year=2026", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/reports/kpis?month=5", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReportsEmptyPeriodReturns200(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/reports/kpis?year=1990&month=1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var kpis models.ReportKPIs
	err := json.Unmarshal(rec.Body.Bytes(), &kpis)
	assert.NoError(t, err)
	assert.Equal(t, 0, kpis.Income)
	assert.Equal(t, 0, kpis.Expense)
	assert.Equal(t, 0, kpis.Balance)
}
