package news

import (
	"context"

	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
)

//go:generate moq -out newsRepo_moq.go . Repository
type Repository interface {
	Save(context.Context, News) error
	GetAll(context.Context) ([]News, error)
	GetById(context.Context, uuid.UUID) (*News, error)
	GetAllByTopic(context.Context, uuid.UUID) ([]News, error)
	GetAllByStatus(ctx context.Context, status bareknews.Status) ([]News, error)
	Count(context.Context, uuid.UUID) (int, error)
	Update(context.Context, News) error
	Delete(context.Context, uuid.UUID) error
}
