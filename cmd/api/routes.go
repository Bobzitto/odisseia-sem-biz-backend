package main

//test
import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (app *application) routes() *echo.Echo {
	e := echo.New()

	SetupMiddleware(e)

	// Use Echo's recover middleware
	e.Use(middleware.Recover())

	// Use CORS middleware (or any other custom middleware)

	// Unprotected route
	e.POST("/authenticate", app.authenticate)
	e.GET("/refresh", app.refreshToken)
	// Protected routes group
	//protected := e.Group("/home")
	//protected.Use(app.AuthRequired)
	//protected.GET("", someProtectedHandler)

	return e
}
