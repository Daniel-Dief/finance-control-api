package middlewares

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func BindMiddlewares(echoAPI *echo.Echo) {
	echoAPI.Use(middleware.RequestLogger())
	echoAPI.Use(middleware.Recover())
	echoAPI.Use(middleware.BodyLimit(1 << 20)) // 1 MB

	CorsMiddleware(echoAPI)
	RatelimitMiddleware(echoAPI)
}
