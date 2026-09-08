package middlewares

import (
	"finance-control-api/database"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		SuccessHandler: func(c *echo.Context) error {
			token, err := echo.ContextGet[*jwt.Token](c, "user")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Token inválido"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Falha ao processar claims do token"})
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Claims do token inválidas"})
			}

			user, err := database.GetUserByID(int(userIDFloat))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Usuário não encontrado"})
			}

			c.Set("user", user)
			return nil
		},
	})
}
