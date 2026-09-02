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

func TestListAreas(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/areas/list", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusNoContent)
}

func TestListAreasWithNameFilter(t *testing.T) {
	e := setupEcho()

	_, err := createTestArea("Test Area for List")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/areas/list?name=Test", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var areas []models.Area
	err = json.Unmarshal(rec.Body.Bytes(), &areas)
	assert.NoError(t, err)
	assert.NotEmpty(t, areas)

	for _, a := range areas {
		database.DeleteArea(a.ID)
	}
}

func TestGetAreaByID(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Area for Get Test")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/areas/%d", area.ID), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var responseArea models.Area
	err = json.Unmarshal(rec.Body.Bytes(), &responseArea)
	assert.NoError(t, err)
	assert.Equal(t, area.ID, responseArea.ID)

	database.DeleteArea(area.ID)
}

func TestGetAreaByIDNotFound(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/areas/99999", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAreaByIDInvalid(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodGet, "/areas/abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateArea(t *testing.T) {
	e := setupEcho()

	body := `{"name":"New Test Area"}`
	req := httptest.NewRequest(http.MethodPost, "/areas/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["id"])
	assert.Equal(t, "New Test Area", response["name"])

	if id, ok := response["id"].(float64); ok {
		database.DeleteArea(int(id))
	}
}

func TestCreateAreaInvalidBody(t *testing.T) {
	e := setupEcho()

	body := `{"invalid": "json"}`
	req := httptest.NewRequest(http.MethodPost, "/areas/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateAreaEmptyName(t *testing.T) {
	e := setupEcho()

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPost, "/areas/create", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpdateArea(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Area to Update")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	body := `{"name":"Updated Area Name"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/areas/%d", area.ID), strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response models.Area
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Area Name", response.Name)

	database.DeleteArea(area.ID)
}

func TestUpdateAreaNotFound(t *testing.T) {
	e := setupEcho()

	body := `{"name":"Updated Name"}`
	req := httptest.NewRequest(http.MethodPut, "/areas/99999", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateAreaInvalidID(t *testing.T) {
	e := setupEcho()

	body := `{"name":"Updated Name"}`
	req := httptest.NewRequest(http.MethodPut, "/areas/abc", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteArea(t *testing.T) {
	e := setupEcho()

	area, err := createTestArea("Area to Delete")
	if err != nil {
		t.Fatalf("Failed to create test area: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/areas/%d", area.ID), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Área excluída com sucesso", response["info"])

	database.DeleteArea(area.ID)
}

func TestDeleteAreaNotFound(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/areas/99999", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteAreaInvalidID(t *testing.T) {
	e := setupEcho()

	req := httptest.NewRequest(http.MethodDelete, "/areas/abc", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}