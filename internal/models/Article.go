package models

import "time"

type Article struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ImageUrl        string    `json:"imageUrl"`
	PublicationDate time.Time `json:"publicationDate"`
	IsPublished     bool      `json:"isPublished"`
	AuthorId        int       `json:"authorId"`
	CategoryId      int       `json:"categoryId"`
	DomainId        int       `json:"domainId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}
