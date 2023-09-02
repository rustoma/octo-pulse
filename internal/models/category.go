package models

import "time"

type Category struct {
	ID           int       `json:"id"`
	CategoryName string    `json:"categoryName"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"UpdatedAt"`
}
