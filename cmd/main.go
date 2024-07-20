package main

import (
	"os"

	"github.com/dmcleish91/go-pet-agency/internal/models"
	"github.com/joho/godotenv"
)

type application struct {
	users *models.UserModel
	pets  *models.PetModel
}

func main() {
	godotenv.Load()
	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn := CreateDatabaseConnection(DATABASE_URL)
	defer conn.Close()

	app := &application{
		users: &models.UserModel{DB: conn},
		pets:  &models.PetModel{DB: conn},
	}

	e := app.Routes()

	e.Logger.Fatal(e.Start(":1323"))
}
