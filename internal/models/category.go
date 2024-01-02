package models

import "time"

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Slug      string    `json:"slug" validate:"required"`
	Weight    int       `json:"weight" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
}
