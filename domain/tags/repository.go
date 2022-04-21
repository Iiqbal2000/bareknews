package tags

import (
	"context"

	"github.com/google/uuid"
)

//go:generate moq -out tagRepo_moq.go . Repository
type Repository interface {
	Save(context.Context, Tags) error
	Update(context.Context, Tags) error
	Delete(context.Context, uuid.UUID) error
	GetById(context.Context, uuid.UUID) (*Tags, error)
	GetAll(context.Context) ([]Tags, error)
	Count(context.Context, uuid.UUID) (int, error)
	GetByNames(context.Context, ...string) ([]Tags, error)
	GetByIds(context.Context, []uuid.UUID) ([]Tags, error)
}
