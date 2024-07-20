package main

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *application) Routes() *echo.Echo {
	signingKey := os.Getenv("SigningKey")

	e := echo.New()

	e.Use(ServerHeader)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339} ${status} ${method} ${host}${path} ${latency_human}]` + "\n",
	}))

	secured := e.Group("/secure")

	secured.Use(echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS256",
		SigningKey:    []byte(signingKey),
	}))

	e.GET("/available", app.GetAllAvailableAnimals)

	e.GET("/details", app.GetPetDetails)

	e.POST("/register", app.Register)

	e.POST("/login", app.Login)

	secured.POST("/addlisting", app.AddAdoptionInformation)

	secured.PUT("/editlisting", app.EditAdoptionInformation)

	secured.PUT("/toggleAdoptionStatus", app.UpdatePetAdoptionStatus)

	return e
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "PetAgency/0.1")

		return next(c)
	}
}
