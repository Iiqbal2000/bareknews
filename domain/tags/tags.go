package tags

import (
	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/google/uuid"
)

// Tags is an aggregrate that represent a tag
type Tags struct {
	Label domain.Label
	Slug  domain.Slug
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
	return t.Label.Validate()
}
