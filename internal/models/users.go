package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}
