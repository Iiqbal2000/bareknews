package news

import "github.com/google/uuid"

//go:generate moq -out newsRepo_moq.go . Repository
type Repository interface {
	Save(News) error
	GetAll() ([]News, error)
	GetById(string) (*News, error)
	// GetByTopic(string) error
	Update(News) error
	Delete(uuid.UUID) error
}