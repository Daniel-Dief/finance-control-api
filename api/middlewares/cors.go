package middlewares

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func CorsMiddleware(echoAPI *echo.Echo) {
	config := middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}

	if os.Getenv("ENV") == "production" {
		config.AllowOrigins = []string{os.Getenv("FRONTEND_URL")}
	}

	echoAPI.Use(middleware.CORSWithConfig(config))
}
