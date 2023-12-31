package models

import "time"

type Article struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Body            string    `json:"body"`
	Thumbnail       *int      `json:"thumbnail"`
	PublicationDate time.Time `json:"publicationDate"`
	IsPublished     bool      `json:"isPublished"`
	AuthorId        int       `json:"authorId"`
	CategoryId      int       `json:"categoryId"`
	DomainId        int       `json:"domainId"`
	Featured        bool      `json:"featured"`
	ReadingTime     *int      `json:"readingTime"`
	IsSponsored     bool      `json:"isSponsored"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
