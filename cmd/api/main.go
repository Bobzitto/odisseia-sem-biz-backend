package main

import (
	"github.com/labstack/echo"
)

func main() {
	//new instance
	e := echo.New()

	//middleware
	SetupMiddleware(e)

	//server on port 8081
	e.Logger.Fatal(e.Start(":8081"))
}
