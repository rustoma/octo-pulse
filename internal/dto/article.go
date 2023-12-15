package dto

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"time"
)

type Article struct {
	ID              int             `json:"id"`
	Title           string          `json:"title" validate:"required,min=4"`
	Slug            string          `json:"slug" validate:"required"`
	Description     string          `json:"description" validate:"required,min=4"`
	Thumbnail       *models.Image   `json:"thumbnail" validate:"required"`
	PublicationDate time.Time       `json:"publicationDate" validate:"required,min=4"`
	IsPublished     bool            `json:"isPublished" validate:"required,min=4"`
	Author          models.Author   `json:"author" validate:"required"`
	Category        models.Category `json:"category" validate:"required"`
	DomainId        int             `json:"domainId" validate:"required,min=4"`
	Featured        bool            `json:"featured"`
	CreatedAt       time.Time       `json:"createdAt" validate:"required,min=4"`
	UpdatedAt       time.Time       `json:"UpdatedAt" validate:"required,min=4"`
}

type ArticleValidationErrors struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Description     string `json:"description"`
	Thumbnail       string `json:"thumbnail"`
	PublicationDate string `json:"publicationDate"`
	IsPublished     string `json:"isPublished"`
	Author          string `json:"author"`
	Category        string `json:"category"`
	DomainId        string `json:"domainId"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"UpdatedAt"`
}
