package middlewares

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func BindMiddlewares(echoAPI *echo.Echo) {
	echoAPI.Use(middleware.RequestLogger())
	echoAPI.Use(middleware.Recover())

	RatelimitMiddleware(echoAPI)
}
