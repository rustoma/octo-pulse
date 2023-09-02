package models

import "time"

type Role struct {
	ID        int       `json:"id"`
	RoleName  string    `json:"roleName"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
