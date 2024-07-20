package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashedpassword"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) RegisterNewUser(user User) (int64, error) {
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3);"

	result, err := m.DB.Exec(context.Background(), query, user.Username, user.Email, user.HashedPassword)
	if err != nil {
		return 0, fmt.Errorf("unable to execute the query. %v", err)
	}

	return result.RowsAffected(), nil
}

func (m *UserModel) GetUserByEmail(email string) (*User, error) {
	query := "SELECT id, username, email, password, created_at FROM users WHERE email = $1"

	row := m.DB.QueryRow(context.Background(), query, email)
	user := &User{}

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("unable to scan the row. %v", err)
	}

	return user, nil
}

func (m *UserModel) EmailExists(email string) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE email = $1"

	var count int
	err := m.DB.QueryRow(context.Background(), query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("unable to execute the query. %v", err)
	}

	return count > 0, nil
}

// these functions can only be used by authenticated users we should get the user email from the JWT Token

func (m *UserModel) UpdateUserEmail(oldEmail string, newEmail string) (int64, error) {
	query := "UPDATE users SET email = $1 WHERE email = $2"

	result, err := m.DB.Exec(context.Background(), query, newEmail, oldEmail)
	if err != nil {
		return 0, fmt.Errorf("unable to execute the query. %v", err)
	}

	return result.RowsAffected(), nil
}

func (m *UserModel) UpdateUsername(email string, newUsername string) (int64, error) {
	query := "UPDATE users SET username = $1 WHERE email = $2"

	result, err := m.DB.Exec(context.Background(), query, newUsername, email)
	if err != nil {
		return 0, fmt.Errorf("unable to execute the query. %v", err)
	}

	return result.RowsAffected(), nil
}
