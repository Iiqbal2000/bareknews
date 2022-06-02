package domain

import (
	"github.com/google/uuid"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Post is an entity that represents a post.
type Post struct {
	ID uuid.UUID
	Title string
	Body string
}

func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.Body, validation.Required),
	)
}