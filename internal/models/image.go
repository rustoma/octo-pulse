package models

import "time"

type Image struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Size       int       `json:"size"`
	Type       string    `json:"type"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Alt        string    `json:"alt"`
	CategoryId int       `json:"categoryId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
