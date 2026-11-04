package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RatelimitMiddleware(echoAPI *echo.Echo) {
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 1 * time.Minute},
		),
		IdentifierExtractor: func(c *echo.Context) (string, error) {
			id := c.RealIP()
			return id, nil
		},
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(c *echo.Context, identifier string, err error) error {
			return c.JSON(http.StatusTooManyRequests, nil)
		},
	}

	echoAPI.Use(middleware.RateLimiterWithConfig(config))
}
