package routes

import (
	"errors"
	"finance-control-api/database"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginBody defines the JSON body accepted by Login.
type LoginBody struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

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
//	@Description	Valida login e senha e retorna um token JWT com o identificador do usuário.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginBody	true	"Credenciais do usuário"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/login [post]
func Login(authGroup *echo.Group) {
	authGroup.POST("/login", func(c *echo.Context) error {
		var body LoginBody

		if !bindAndValidate(c, &body) {
			return nil
		}

		user, err := database.GetUserByLogin(body.Login)
		if errors.Is(err, database.ErrNotFound) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Login ou senha inválidos"})
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
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
