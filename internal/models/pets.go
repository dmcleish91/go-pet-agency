package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pet struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Species           string    `json:"species"`
	Breed             string    `json:"breed"`
	Age               int       `json:"age"`
	Description       string    `json:"description"`
	Status            string    `json:"status"`
	Gender            string    `json:"gender"`
	Size              string    `json:"size"`
	Color             string    `json:"color"`
	Weight            float32   `json:"weight"`
	VaccinationStatus bool      `json:"vaccinationstatus"`
	Spayed            bool      `json:"spayed"`
	Microchipped      bool      `json:"microchipped"`
	RescueStory       string    `json:"rescuestory"`
	CreatedAt         time.Time `json:"createdat"`
	UserID            int       `json:"userid"`
}

type PetModel struct {
	DB *pgxpool.Pool
}

func (m *PetModel) GetAllAvailablePets() ([]*Pet, error) {
	query := `SELECT id, name, species, breed, age, status, gender, size, color, weight, vaccination_status, spayed, microchipped, 
	rescue_story, created_at, user_id FROM pets WHERE status = 'available'`

	rows, err := m.DB.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("unable to execute the query. %v", err)
	}
	defer rows.Close()

	var pets []*Pet
	for rows.Next() {
		pet := &Pet{}
		err := rows.Scan(&pet.ID, &pet.Name, &pet.Species, &pet.Breed, &pet.Age, &pet.Status, &pet.Gender, &pet.Size, &pet.Color,
			&pet.Weight, &pet.VaccinationStatus, &pet.Spayed, &pet.Microchipped, &pet.RescueStory, &pet.CreatedAt, &pet.UserID)
		if err != nil {
			return nil, fmt.Errorf("unable to scan the row. %v", err)
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

func (m *PetModel) GetPetDetails(petID int) (*Pet, error) {
	query := `SELECT id, name, species, breed, age, description, status, gender, size, color, weight, vaccination_status, spayed, 
	microchipped, rescue_story, created_at FROM pets WHERE id = $1`

	row := m.DB.QueryRow(context.Background(), query, petID)
	pet := &Pet{}

	err := row.Scan(&pet.ID, &pet.Name, &pet.Species, &pet.Breed, &pet.Age, &pet.Description, &pet.Status, &pet.Gender, &pet.Size,
		&pet.Color, &pet.Weight, &pet.VaccinationStatus, &pet.Spayed, &pet.Microchipped, &pet.RescueStory, &pet.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to scan the row. %v", err)
	}

	return pet, nil
}

func (m *PetModel) AddPetListing(pet Pet) (int64, error) {
	query := `INSERT INTO pets (name, species, breed, age, description, status, gender, size, color, weight, vaccination_status, 
	spayed, microchipped, rescue_story, user_id) VALUES ($1, $2, $3, $4, $5, 'available', $6, $7, $8, $9, $10, $11, $12, $13, $14);`

	result, err := m.DB.Exec(context.Background(), query, pet.Name, pet.Species, pet.Breed, pet.Age, pet.Description, pet.Gender,
		pet.Size, pet.Color, pet.Weight, pet.VaccinationStatus, pet.Spayed, pet.Microchipped, pet.RescueStory, pet.UserID)
	if err != nil {
		return 0, fmt.Errorf("unable to execute the query. %v", err)
	}

	return result.RowsAffected(), nil
}

func (m *PetModel) EditPetListing(pet Pet) (int64, error) {
	query := `UPDATE pets SET name = $1, species = $2, breed = $3, age = $4, description = $5, gender = $6, size = $7, color = $8, 
	weight = $9, vaccination_status = $10, spayed = $11, microchipped = $12, rescue_story = $13 WHERE id = $14 AND user_id = $15`

	result, err := m.DB.Exec(context.Background(), query, pet.Name, pet.Species, pet.Breed, pet.Age, pet.Description, pet.Gender,
		pet.Size, pet.Color, pet.Weight, pet.VaccinationStatus, pet.Spayed, pet.Microchipped, pet.RescueStory, pet.ID, pet.UserID)
	if err != nil {
		return 0, fmt.Errorf("unable to execute the query. %v", err)
	}

	return result.RowsAffected(), nil
}

func (m *PetModel) TogglePetStatus(petID int, userID int, status string) (int64, error) {
	query := "UPDATE pets SET status = $1 WHERE id = $2 AND user_id = $3"

	result, err := m.DB.Exec(context.Background(), query, status, petID, userID)
	if err != nil {
		return 0, fmt.Errorf("unable to execute the query. %v", err)
	}

	return result.RowsAffected(), nil
}
