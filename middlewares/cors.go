package middlewares

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				"GET",
				"POST",
				"OPTIONS",
			},
			AllowHeaders: []string{
				"Content-Type",
				"Authorization",
			},
		},
	)
}
