package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"finance-control-api/database"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

// TestAreaNameUniquePerUser verifies that two different users can own an area
// with the same name.
func TestAreaNameUniquePerUser(t *testing.T) {
	other, err := database.CreateUser("Outro Usuário", "outrousuario", "hash")
	if err != nil {
		t.Fatalf("Failed to create second user: %v", err)
	}

	areaA, err := database.CreateArea(testUserID, "teste")
	assert.NoError(t, err)
	areaB, err := database.CreateArea(other.ID, "teste")
	assert.NoError(t, err)
	assert.NotEqual(t, areaA.ID, areaB.ID)
	assert.NotEqual(t, areaA.UserID, areaB.UserID)

	assert.NoError(t, database.DeleteArea(testUserID, areaA.ID))
	assert.NoError(t, database.DeleteArea(other.ID, areaB.ID))
}

// TestCategoryNameUniquePerUser verifies that two different users can own a
// category with the same name.
func TestCategoryNameUniquePerUser(t *testing.T) {
	other, err := database.CreateUser("Outro Usuário 2", "outrousuario2", "hash")
	if err != nil {
		t.Fatalf("Failed to create second user: %v", err)
	}

	catA, err := database.CreateCategory(testUserID, "teste")
	assert.NoError(t, err)
	catB, err := database.CreateCategory(other.ID, "teste")
	assert.NoError(t, err)
	assert.NotEqual(t, catA.ID, catB.ID)
	assert.NotEqual(t, catA.UserID, catB.UserID)

	assert.NoError(t, database.DeleteCategory(testUserID, catA.ID))
	assert.NoError(t, database.DeleteCategory(other.ID, catB.ID))
}

// TestCreateAreaDuplicateNameConflict verifies that the same user cannot have
// two active areas with the same name (HTTP 409).
func TestCreateAreaDuplicateNameConflict(t *testing.T) {
	e := setupEcho()

	post := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/areas/create", strings.NewReader(`{"name":"nome_duplicado"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		return rec
	}

	rec := post()
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to parse created area: %v", err)
	}

	rec2 := post()
	assert.Equal(t, http.StatusConflict, rec2.Code)

	assert.NoError(t, database.DeleteArea(testUserID, created.ID))
}

// TestAreaNameReusableAfterSoftDelete verifies that the name of a soft-deleted
// area can be reused.
func TestAreaNameReusableAfterSoftDelete(t *testing.T) {
	e := setupEcho()

	create := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/areas/create", strings.NewReader(`{"name":"reutilizavel"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		return rec
	}

	rec := create()
	assert.Equal(t, http.StatusCreated, rec.Code)

	var created struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to parse created area: %v", err)
	}

	delReq := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/areas/%d", created.ID), nil)
	delRec := httptest.NewRecorder()
	e.ServeHTTP(delRec, delReq)
	assert.Equal(t, http.StatusOK, delRec.Code)

	rec2 := create()
	assert.Equal(t, http.StatusCreated, rec2.Code)

	var recreated struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(rec2.Body.Bytes(), &recreated); err != nil {
		t.Fatalf("Failed to parse recreated area: %v", err)
	}
	assert.NoError(t, database.DeleteArea(testUserID, recreated.ID))
}

// TestBudgetReusableAfterSoftDelete verifies that a budget for the same
// area/period can be recreated after the previous one was soft-deleted.
func TestBudgetReusableAfterSoftDelete(t *testing.T) {
	area, err := database.CreateArea(testUserID, "area_orcamento")
	if err != nil {
		t.Fatalf("Failed to create area: %v", err)
	}
	defer database.DeleteArea(testUserID, area.ID)

	year, month, amount := 2026, 1, 1000
	props := database.BudgetProps{
		Year:   &year,
		Month:  &month,
		AreaId: &area.ID,
		Amount: &amount,
	}

	b1, err := database.CreateBudget(testUserID, props)
	assert.NoError(t, err)

	assert.NoError(t, database.DeleteBudget(testUserID, b1.ID))

	_, err = database.CreateBudget(testUserID, props)
	assert.NoError(t, err)
}