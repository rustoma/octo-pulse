package models

import "time"

type BasicPage struct {
	ID        int       `json:"id"`
	Title     string    `json:"title" validate:"required"`
	Slug      string    `json:"slug" validate:"required"`
	Body      string    `json:"body" validate:"required"`
	Domain    int       `json:"domain" validate:"required"`
	CreatedAt time.Time `json:"createdAt" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
}
