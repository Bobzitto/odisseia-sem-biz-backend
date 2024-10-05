package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func registerRoutes(e *echo.Echo) {
	e.GET("/", loginHandler)
	e.GET("/home", homeHandler)
	e.GET("/home/aulas", aulasHandler)
	e.GET("/home/turmas", turmasHandler)
	e.POST("/home/turmas/0/edit", criaTurmaHandler)
	e.POST("/home/aulas/0/edit", criaAulaHandler)

}

func homeHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome")
}
