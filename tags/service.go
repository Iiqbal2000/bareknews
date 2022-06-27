package tags

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
)

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
	tag, err := s.store.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TagsOut{}, bareknews.ErrDataNotFound
		} else {
			return TagsOut{}, err
		}
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
	_, err := s.store.Count(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return bareknews.ErrDataNotFound
		} else {
			return err
		}
	}

	err = s.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetById(ctx context.Context, id uuid.UUID) (TagsOut, error) {
	tg, err := s.store.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TagsOut{}, bareknews.ErrDataNotFound
		} else {
			return TagsOut{}, err
		}
	}

	return TagsOut{
		ID:   tg.Label.ID,
		Name: tg.Label.Name,
		Slug: tg.Slug.String(),
	}, nil
}

func (s Service) GetByIds(ctx context.Context, ids []uuid.UUID) ([]TagsOut, error) {
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
