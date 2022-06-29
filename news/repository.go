package news

import (
	"context"

	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
)

//go:generate moq -out newsRepo_moq.go . Repository
type Repository interface {
	Save(context.Context, News) error
	GetAll(ctx context.Context, cursor int64, limit int) ([]News, error)
	GetById(context.Context, uuid.UUID) (*News, error)
	GetAllByTopic(ctx context.Context, id uuid.UUID, cursor int64, limit int) ([]News, error)
	GetAllByStatus(ctx context.Context, status bareknews.Status, cursor int64, limit int) ([]News, error)
	Count(context.Context, uuid.UUID) (int, error)
	Update(context.Context, News) error
	Delete(context.Context, uuid.UUID) error
}
