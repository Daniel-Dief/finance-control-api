package routes

import (
	"encoding/base64"
	"errors"
	"finance-control-api/database"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterBody defines the JSON body accepted by Register.
type RegisterBody struct {
	Name     string `json:"name" validate:"required"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginResponse defines the response returned by a successful Login.
type LoginResponse struct {
	Token string      `json:"token"`
	User  UserPayload `json:"user"`
}

type UserPayload struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

// Login godoc
//
//	@Summary		Autenticar usuário
//	@Description	Valida as credenciais enviadas no header Authorization, no formato Basic base64(usuario:senha), e retorna um token JWT com o identificador do usuário.
//	@Tags			auth
//	@Produce		json
//	@Param			authorization	header		string	true	"Credenciais no formato Basic base64(usuario:senha)"
//	@Success		200				{object}	LoginResponse
//	@Failure		400				{object}	map[string]interface{}
//	@Failure		401				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/auth/login [post]
func Login(authGroup *echo.Group) {
	authGroup.POST("/login", func(c *echo.Context) error {
		login, password, ok := basicAuthCredentials(c)
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Credenciais inválidas no cabeçalho Authorization"})
		}

		user, err := database.GetUserByLogin(login)
		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Login ou senha inválidos"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Login ou senha inválidos"})
		}

		claims := jwt.MapClaims{
			"user_id": user.ID,
			"sub":     user.Login,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
			"iat":     time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Falha ao gerar o token de autenticação"})
		}

		return c.JSON(http.StatusOK, LoginResponse{
			Token: tokenString,
			User: UserPayload{
				ID:    user.ID,
				Name:  user.Name,
				Login: user.Login,
			},
		})
	})
}

// basicAuthCredentials extracts the login and password from the Authorization
// header. The value must contain "usuario:senha" encoded in base64 and may be
// prefixed with "Basic ". It returns false when the header is missing or
// malformed.
func basicAuthCredentials(c *echo.Context) (string, string, bool) {
	header := strings.TrimSpace(c.Request().Header.Get(echo.HeaderAuthorization))
	if header == "" {
		return "", "", false
	}

	// Accept both the standard "Basic <credentials>" scheme and the raw
	// base64 value.
	if len(header) >= len("Basic ") && strings.EqualFold(header[:len("Basic ")], "Basic ") {
		header = strings.TrimSpace(header[len("Basic "):])
	}

	decoded, err := base64.StdEncoding.DecodeString(header)
	if err != nil {
		return "", "", false
	}

	login, password, ok := strings.Cut(string(decoded), ":")
	if !ok || login == "" || password == "" {
		return "", "", false
	}

	return login, password, true
}

// Register godoc
//
//	@Summary		Criar usuário
//	@Description	Cria um novo usuário com status ativo. O login deve ser único e a senha precisa de no mínimo 6 caracteres.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body	RegisterBody	true	"Dados do novo usuário"
//	@Success		201		{object}	UserPayload
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/register [post]
func Register(authGroup *echo.Group) {
	authGroup.POST("/register", func(c *echo.Context) error {
		var body RegisterBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Falha ao processar a senha"})
		}

		user, err := database.CreateUser(body.Name, body.Login, string(hash))
		if errors.Is(err, database.ErrDuplicateLogin) {
			return c.JSON(http.StatusConflict, map[string]string{"error": database.ErrDuplicateLogin.Error()})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusCreated, UserPayload{
			ID:    user.ID,
			Name:  user.Name,
			Login: user.Login,
		})
	})
}
