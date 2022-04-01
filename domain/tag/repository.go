package tag

import "github.com/Iiqbal2000/bareknews"

//go:generate moq -out tagRepo_moq.go . Repository
type Repository interface {
	Save(bareknews.Tags) error
	Update(bareknews.Tags) error
	GetById(string) (*bareknews.Tags, error)
	GetByNewsId(id string) ([]bareknews.Tags, error)
	GetByName(names string) (bareknews.Tags, error)
	GetByNames(names ...string) ([]bareknews.Tags, error)
	GetAll() ([]bareknews.Tags, error)
	Delete(id string) error
}
