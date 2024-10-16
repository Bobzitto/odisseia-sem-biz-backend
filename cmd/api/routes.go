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
	e.GET("/materias", app.TodasMaterias)
	e.POST("/graph", app.AulasGraphQL)

	// Protected routes group
	protected := e.Group("/home")
	protected.Use(app.AuthRequired)
	protected.GET("/aulas", app.TodasAulas)
	protected.PUT("/aulas/0", app.InserirAula)
	protected.PATCH("/aulas/:id", app.AtualizarAula)
	protected.GET("/aulas/{id}", app.EditarAula)
	protected.DELETE("/aulas/:id", app.DeletarAula)

	//turmas

	protected.GET("/turmas", app.TodasTurmas)
	protected.PUT("/turmas/0", app.InserirTurma)
	protected.PATCH("/turmas/:id", app.AtualizarTurma)
	protected.GET("/turmas/:id", app.EditarTurma)
	protected.DELETE("/turmas/:id", app.DeletarTurma)

	return e
}
