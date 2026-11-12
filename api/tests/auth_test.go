package tests

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"finance-control-api/routes"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

// basicAuthHeader builds the Authorization header value from the given
// credentials encoded as base64(usuario:senha).
func basicAuthHeader(login, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password))
}

func TestLoginSuccess(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set(echo.HeaderAuthorization, basicAuthHeader(testUserLogin, testUserPassword))
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

func TestLoginSuccessWithRawBase64(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	credentials := base64.StdEncoding.EncodeToString([]byte(testUserLogin + ":" + testUserPassword))
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set(echo.HeaderAuthorization, credentials)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestLoginWrongPassword(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set(echo.HeaderAuthorization, basicAuthHeader(testUserLogin, "wrong"))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginUnknownUser(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set(echo.HeaderAuthorization, basicAuthHeader("unknown", "password"))
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginMissingAuthorization(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginMalformedAuthorization(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set(echo.HeaderAuthorization, "Basic not-base64")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginMissingPassword(t *testing.T) {
	e := echo.New()
	routes.BindRoutes(e)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	req.Header.Set(echo.HeaderAuthorization, basicAuthHeader(testUserLogin, ""))
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
