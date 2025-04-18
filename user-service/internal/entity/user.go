package models

import "time"

type User struct {
	ID           string `db:"id"`
	Email        string `db:"email"`
	Password     string `db:"password_hash"`
	Role         string `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}