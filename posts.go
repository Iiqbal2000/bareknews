package bareknews

import (
	"github.com/google/uuid"
)

// Post is an entity that represents a post
type Posts struct {
	ID uuid.UUID
	Title string
	Status Status
	Body string
}