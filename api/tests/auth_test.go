package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"finance-control-api/routes"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
	e := setupEcho()

	body := fmt.Sprintf(`{"login":"%s","password":"%s"}`, testUserLogin, testUserPassword)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response struct {
		Token string `json:"token"`
		User  struct {
			ID    int    `json:"id"`
			Login string `json:"login"`
		} `json:"user"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, testUserID, response.User.ID)
	assert.Equal(t, testUserLogin, response.User.Login)
}

func TestLoginWrongPassword(t *testing.T) {
	e := setupEcho()

	body := fmt.Sprintf(`{"login":"%s","password":"wrong"}`, testUserLogin)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginUnknownUser(t *testing.T) {
	e := setupEcho()

	body := `{"login":"unknown","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginMissingFields(t *testing.T) {
	e := setupEcho()

	body := `{"login":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProtectedRouteWithoutToken(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	req := httptest.NewRequest(http.MethodGet, "/areas/list", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestProtectedRouteWithInvalidToken(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	req := httptest.NewRequest(http.MethodGet, "/areas/list", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer invalid.token.here")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRegisterSuccess(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	body := `{"name":"Novo Usuário","login":"novouser","password":"senha123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Login string `json:"login"`
	}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotZero(t, response.ID)
	assert.Equal(t, "Novo Usuário", response.Name)
	assert.Equal(t, "novouser", response.Login)
}

func TestRegisterDuplicateLogin(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	body := `{"name":"Duplicado","login":"novouser","password":"senha123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestRegisterMissingFields(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	body := `{"name":"","login":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterWeakPassword(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	body := `{"name":"Fulano","login":"fulano","password":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}