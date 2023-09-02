package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}
