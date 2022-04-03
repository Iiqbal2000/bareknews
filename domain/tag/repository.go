package tag

import (
	"github.com/Iiqbal2000/bareknews/domain"
)

//go:generate moq -out tagRepo_moq.go . Repository
type Repository interface {
	Save(domain.Tags) error
	Update(domain.Tags) error
	GetById(string) (*domain.Tags, error)
	GetByNewsId(id string) ([]domain.Tags, error)
	GetByName(names string) (domain.Tags, error)
	GetByNames(names ...string) ([]domain.Tags, error)
	GetAll() ([]domain.Tags, error)
	Delete(id string) error
}
