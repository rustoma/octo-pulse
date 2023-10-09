package dto

import "time"

type Article struct {
	ID              int       `json:"id"`
	Title           string    `json:"title" validate:"required,min=4"`
	Description     string    `json:"description" validate:"required,min=4"`
	ImageUrl        string    `json:"imageUrl" validate:"required,min=4"`
	PublicationDate time.Time `json:"publicationDate" validate:"required,min=4"`
	IsPublished     bool      `json:"isPublished" validate:"required,min=4"`
	AuthorId        int       `json:"authorId" validate:"required,min=4"`
	CategoryId      int       `json:"categoryId" validate:"required,min=4"`
	DomainId        int       `json:"domainId" validate:"required,min=4"`
	CreatedAt       time.Time `json:"createdAt" validate:"required,min=4"`
	UpdatedAt       time.Time `json:"UpdatedAt" validate:"required,min=4"`
}

type ArticleValidationErrors struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	ImageUrl        string `json:"imageUrl"`
	PublicationDate string `json:"publicationDate"`
	IsPublished     string `json:"isPublished"`
	AuthorId        string `json:"authorId"`
	CategoryId      string `json:"categoryId"`
	DomainId        string `json:"domainId"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"UpdatedAt"`
}
