package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dmcleish91/go-pet-agency/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	Sub   string `json:"sub"`  // Subject (user ID)
	Name  string `json:"name"` // Name of the user
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func (app *application) GetAllAvailableAnimals(c echo.Context) error {

	allPets, err := app.pets.GetAllAvailablePets()
	if err != nil {
		log.Printf("Error fetching Data: %s", err)
		return c.String(http.StatusInternalServerError, "Error fetching available information")
	}

	return c.JSON(http.StatusOK, allPets)
}

func (app *application) GetPetDetails(c echo.Context) error {
	petId := c.QueryParam("petId")

	id, err := strconv.ParseInt(petId, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %s", err)
		return c.String(http.StatusUnprocessableEntity, "You need to provide a valid number")
	}

	pet, err := app.pets.GetPetDetails(int(id))
	if err != nil {
		log.Printf("Error getting pet details: %s", err)
		return c.String(http.StatusInternalServerError, "Something went wrong. Please try again")
	}

	return c.JSON(http.StatusOK, pet)
}

func (app *application) AddAdoptionInformation(c echo.Context) error {
	pet := models.Pet{}

	defer c.Request().Body.Close()
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &pet)
	if err != nil {
		log.Printf("Failed unmarshaling: %s", err)
		return c.String(http.StatusInternalServerError, "Parsing Error! Is this JSON?")
	}

	userId := GetUserIdFromToken(c)
	if userId == -1 {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	pet.UserID = userId

	result, err := app.pets.AddPetListing(pet)
	if err != nil {
		log.Printf("Postgresql Error: %s", err)
		return c.String(http.StatusInternalServerError, "Database Error: "+err.Error())
	}

	log.Printf("this is your data: %#v", result)
	return c.String(http.StatusOK, "Successfuly added adoption information")
}

func (app *application) EditAdoptionInformation(c echo.Context) error {
	pet := models.Pet{}

	defer c.Request().Body.Close()
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &pet)
	if err != nil {
		log.Printf("Failed unmarshaling: %s", err)
		return c.String(http.StatusInternalServerError, "Parsing Error! Is this JSON?")
	}

	userId := GetUserIdFromToken(c)
	if userId == -1 {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	pet.UserID = userId

	result, err := app.pets.EditPetListing(pet)
	if err != nil {
		log.Printf("Postgresql Error: %s", err)
		return c.String(http.StatusInternalServerError, "Something went wrong. Try again later")
	}

	if result == 0 {
		return c.String(http.StatusUnprocessableEntity, "Error updating adoption information. Did you include the \"id\"")
	}

	log.Printf("User with id %#v modified %#v rows", userId, result)
	return c.String(http.StatusOK, "Successfuly edited adoption information")
}

func (app *application) UpdatePetAdoptionStatus(c echo.Context) error {
	petIdStr := c.QueryParam("petId")

	petId, err := strconv.ParseInt(petIdStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing id: %s", err)
		return c.String(http.StatusUnprocessableEntity, "You need to provide a valid number")
	}

	userId := GetUserIdFromToken(c)
	if userId == -1 {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	pet, err := app.pets.GetPetDetails(int(petId))
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.String(http.StatusInternalServerError, "Something went wrong. Try again later")
	}

	adoptionStatus := pet.Status

	var newAdoptionStatus string
	if adoptionStatus == "available" {
		newAdoptionStatus = "adopted"
	} else {
		newAdoptionStatus = "available"
	}

	result, err := app.pets.TogglePetStatus(int(petId), userId, newAdoptionStatus)
	if err != nil {
		log.Printf("Database Error: %s", err)
		return c.String(http.StatusInternalServerError, "Something went wrong. Try again later")
	}

	if result == 0 {
		return c.String(http.StatusUnprocessableEntity, "Error updating adoption information. Did you include the \"id\"")
	}

	log.Printf("User with id %#v modified %#v rows with new status: %s", userId, result, newAdoptionStatus)
	return c.String(http.StatusOK, "Adoption status set to: "+newAdoptionStatus)
}

func (app *application) Register(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" || password == "" || username == "" {
		return c.String(http.StatusBadRequest, "Email, username and password are required")
	}

	existingUser, err := app.users.EmailExists(email)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong try")
	}

	if existingUser {
		return c.String(http.StatusConflict, "User with supplied email already exists")
	}

	// TODO: Password Strength Function

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong try")
	}

	user := models.User{
		Username:       username,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	newUser, err := app.users.RegisterNewUser(user)
	if err != nil {
		log.Printf("Postgresql Error: %s", err)
		return c.String(http.StatusInternalServerError, "Error registering new user")
	}

	log.Printf("Created %#v new row", newUser)
	return c.String(http.StatusOK, "Successfully registered user")
}

func (app *application) Login(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	user, err := app.users.GetUserByEmail(email)
	if err != nil {
		log.Printf("Error: %s", err)
		return ctx.String(http.StatusInternalServerError, "Something went wrong")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return ctx.String(http.StatusUnauthorized, "email or password was incorrect")
	}

	if subtle.ConstantTimeCompare([]byte(email), []byte(user.Email)) != 1 {
		return ctx.String(http.StatusUnauthorized, "email or password was incorrect")
	}

	token, err := createJwtToken(user.ID, user.Username)
	if err != nil {
		log.Println("Error Creating JWT token")
		return ctx.String(http.StatusInternalServerError, "Something went wrong")
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "You were logged in!",
		"token":   token,
	})
}

// These are utility functions

func createJwtToken(userID int, username string) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		Sub:   fmt.Sprintf("%d", userID),
		Name:  username,
		Admin: false,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("itsasecret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserIdFromToken(c echo.Context) int {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*jwtCustomClaims)
	userIdStr := claims.Sub

	log.Println("Decoded user id is", userIdStr)

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Println("Error converting user ID to int64:", err)
		return -1
	}

	return int(userId)
}
