package models

import "time"

type Domain struct {
	ID         int `json:"id"`
	DomainName string
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"UpdatedAt"`
}
