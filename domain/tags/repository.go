package tags

import "github.com/google/uuid"

//go:generate moq -out tagRepo_moq.go . Repository
type Repository interface {
	Save(Tags) error
	Update(Tags) error
	Delete(id uuid.UUID) error
	GetById(id uuid.UUID) (*Tags, error)
	GetAll() ([]Tags, error)
	// GetByNewsId(id string) ([]Tags, error)
	// GetByName(names string) (Tags, error)
	// GetByNames(names ...string) ([]Tags, error)
}
