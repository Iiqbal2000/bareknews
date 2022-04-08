package news

import "github.com/google/uuid"

//go:generate moq -out newsRepo_moq.go . Repository
type Repository interface {
	Save(News) error
	GetAll() ([]News, error)
	GetById(id uuid.UUID) (*News, error)
	Count(id uuid.UUID) (int, error)
	Update(News) error
	Delete(uuid.UUID) error
}