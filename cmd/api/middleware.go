package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func SetupMiddleware(e *echo.Echo){
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://example.com", "http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders: []string{"Authorization"},
		AllowCredentials: true,
	}))
}

