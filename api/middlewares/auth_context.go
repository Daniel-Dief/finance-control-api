package middlewares

import (
	"errors"
	"finance-control-api/models"

	"github.com/labstack/echo/v5"
)

// ErrUserMissing indicates the authenticated user is not present in the context.
var ErrUserMissing = errors.New("usuário não autenticado na requisição")

// GetCurrentUser retrieves the authenticated user stored in the context by
// the JWT middleware. It returns ErrUserMissing when the request does not
// carry an authenticated user.
func GetCurrentUser(c *echo.Context) (models.User, error) {
	user, err := echo.ContextGet[models.User](c, "user")
	if err != nil {
		return models.User{}, ErrUserMissing
	}
	return user, nil
}
