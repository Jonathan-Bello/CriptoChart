package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/Jonathan-Bello/CriptoChart/handler"
)

// Chart cointains the routes for handlers about Charts
func Chart(e *echo.Echo) {
	chart := e.Group("/v1/charts")
	chart.GET("/:currency/:startdate", handler.CreateChart)
	chart.GET("/:currency/:startdate/:enddate", handler.CreateChart)
}
