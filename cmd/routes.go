package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) Routes() {
	app.routes.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
}
