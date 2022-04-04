package domain

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

var (
	ErrBlankTag = errors.New("invalid tag")
)

// Value object
type Tags struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug Slug      `json:"slug"`
}

func NewTags(tagName string) (*Tags, error) {
	tag := Tags{
		ID:   uuid.New(),
		Name: tagName,
		Slug: NewSlug(tagName),
	}

	return &tag, nil
}

func (t *Tags) ChangeName(newName string) {
	t.Name = newName
	t.Slug = NewSlug(t.Name)
}

func (t Tags) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Name, validation.Required, validation.Length(1, 20)),
	)
}
