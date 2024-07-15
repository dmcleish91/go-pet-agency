package main

import (
	"os"

	"github.com/dmcleish91/go-pet-agency/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var DATABASE_URL string
var conn *pgxpool.Pool

type application struct {
	users  *models.UserModel
	routes *echo.Echo
}

func main() {
	godotenv.Load()
	DATABASE_URL = os.Getenv("DATABASE_URL")
	signingKey := os.Getenv("SigningKey")

	conn = CreateDatabaseConnection(DATABASE_URL)
	defer conn.Close()

	e := echo.New()

	app := &application{
		users:  &models.UserModel{DB: conn},
		routes: e,
	}

	app.Routes()

	e.Use(ServerHeader)
	secured := e.Group("/secure")

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339} ${status} ${method} ${host}${path} ${latency_human}]` + "\n",
	}))

	secured.Use(echojwt.WithConfig(echojwt.Config{
		SigningMethod: "HS512",
		SigningKey:    []byte(signingKey),
	}))

	e.Logger.Fatal(e.Start(":1323"))
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "PetAgency/0.1")

		return next(c)
	}
}
