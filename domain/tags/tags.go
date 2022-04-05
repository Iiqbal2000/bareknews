package tags

import (
	"errors"

	"github.com/Iiqbal2000/bareknews/domain"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

var (
	ErrBlankTag = errors.New("invalid tag")
)

// Agregates
type Tags struct {
	Label domain.Label
	Slug  domain.Slug `json:"slug"`
}

func New(tagName string) *Tags {
	tag := Tags{
		Label: domain.Label{
			ID:   uuid.New(),
			Name: tagName,
		},
		Slug: domain.NewSlug(tagName),
	}

	return &tag
}

func (t *Tags) ChangeName(newName string) {
	t.Label.Name = newName
	t.Slug = domain.NewSlug(newName)
}

func (t Tags) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Label),
	)
}
