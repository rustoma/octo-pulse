package models

import "time"

type Author struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"firstName" validate:"required"`
	LastName    string    `json:"lastName" validate:"required"`
	Description string    `json:"description" validate:"required"`
	ImageUrl    string    `json:"imageUrl" validate:"required"`
	CreatedAt   time.Time `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time `json:"updatedAt" validate:"required"`
}
