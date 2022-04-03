package domain

import (
	"github.com/google/uuid"
)

// Post is an entity that represents a post
type Posts struct {
	ID uuid.UUID	`json:"id"`
	Title string	`json:"title"`
	Status Status	`json:"status"`
	Body string	`json:"body"`
}