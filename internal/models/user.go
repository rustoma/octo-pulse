package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	RefreshToken string    `json:"refreshToken"`
	PasswordHash string    `json:"-"`
	RoleId       int       `json:"roleID"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type UserRoles struct {
	Admin  int
	Editor int
}
