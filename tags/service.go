package tags

import (
	"context"
	"strings"

	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/Iiqbal2000/bareknews/tags")

type TagsOut struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type Service struct {
	store Repository
}

func CreateSvc(repo Repository) Service {
	return Service{repo}
}

func (s Service) Create(ctx context.Context, tagName string) (TagsOut, error) {
	ctx, span := tracer.Start(ctx, "tags.Create")
	defer span.End()

	tag := Create(strings.TrimSpace(tagName))

	err := tag.Validate()
	if err != nil {
		return TagsOut{}, err
	}

	err = s.store.Save(ctx, *tag)
	if err != nil {
		return TagsOut{}, err
	}

	return TagsOut{
		ID:   tag.Label.ID,
		Name: tag.Label.Name,
		Slug: tag.Slug.String(),
	}, nil
}

func (s Service) Update(ctx context.Context, id uuid.UUID, newTagname string) (TagsOut, error) {
	ctx, span := tracer.Start(ctx, "tags.Update")
	defer span.End()

	tag, err := s.store.GetById(ctx, id)
	if err != nil {
		return TagsOut{}, err
	}

	tag.ChangeName(strings.TrimSpace(newTagname))

	err = tag.Validate()
	if err != nil {
		return TagsOut{}, err
	}

	err = s.store.Update(ctx, *tag)
	if err != nil {
		return TagsOut{}, err
	}

	return TagsOut{
		ID:   tag.Label.ID,
		Name: tag.Label.Name,
		Slug: tag.Slug.String(),
	}, nil
}

func (s Service) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "tags.Delete")
	defer span.End()

	_, err := s.store.Count(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetById(ctx context.Context, id uuid.UUID) (TagsOut, error) {
	ctx, span := tracer.Start(ctx, "tags.GetById")
	defer span.End()

	tg, err := s.store.GetById(ctx, id)
	if err != nil {
		return TagsOut{}, err
	}

	return TagsOut{
		ID:   tg.Label.ID,
		Name: tg.Label.Name,
		Slug: tg.Slug.String(),
	}, nil
}

func (s Service) GetByIds(ctx context.Context, ids []uuid.UUID) ([]TagsOut, error) {
	ctx, span := tracer.Start(ctx, "tags.GetByIds")
	defer span.End()

	tgs, err := s.store.GetByIds(ctx, ids)
	if err != nil {
		return []TagsOut{}, bareknews.ErrInternalServer
	}

	r := make([]TagsOut, 0)

	for _, t := range tgs {
		r = append(r, TagsOut{
			ID:   t.Label.ID,
			Name: t.Label.Name,
			Slug: t.Slug.String(),
		})
	}

	return r, nil
}

func (s Service) GetAll(ctx context.Context) ([]TagsOut, error) {
	ctx, span := tracer.Start(ctx, "tags.GetAll")
	defer span.End()

	tg, err := s.store.GetAll(ctx)
	if err != nil {
		return []TagsOut{}, err
	}

	r := make([]TagsOut, 0)

	for _, t := range tg {
		r = append(r, TagsOut{
			ID:   t.Label.ID,
			Name: t.Label.Name,
			Slug: t.Slug.String(),
		})
	}

	return r, nil
}

func (s Service) GetByNames(ctx context.Context, names []string) []TagsOut {
	ctx, span := tracer.Start(ctx, "tags.GetByNames")
	defer span.End()

	tg, err := s.store.GetByNames(ctx, names...)
	if err != nil {
		return []TagsOut{}
	}

	r := make([]TagsOut, 0)

	for _, t := range tg {
		r = append(r, TagsOut{
			ID:   t.Label.ID,
			Name: t.Label.Name,
			Slug: t.Slug.String(),
		})
	}

	return r
}

func (s Service) GetByName(ctx context.Context, name string) TagsOut {
	ctx, span := tracer.Start(ctx, "tags.GetByName")
	defer span.End()

	tg, err := s.store.GetByName(ctx, name)
	if err != nil {
		return TagsOut{}
	}

	return TagsOut{
		ID:   tg.Label.ID,
		Name: tg.Label.Name,
		Slug: tg.Slug.String(),
	}
}
