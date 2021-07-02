package main

import (
	"net/http"
	"os"

	"github.com/Jonathan-Bello/CriptoChart/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	routes.Chart(e)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"saludo": "holiwis uwu",
		})
	})

	e.Start(":" + port)
}
