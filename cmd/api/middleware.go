package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetupMiddleware(e *echo.Echo) {
	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000" // Default for local development
	}

	// Split the origins into a slice
	origins := strings.Split(allowOrigins, ",")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     origins,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))
}

func (app *application) AuthRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, _, err := app.auth.GetTokenFromHeaderAndVerify(c)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}
