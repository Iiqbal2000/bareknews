package tags

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
)

type Response struct {
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

func (s Service) Create(ctx context.Context, tagName string) (Response, error) {
	tag := Create(strings.TrimSpace(tagName))

	err := tag.Validate()
	if err != nil {
		return Response{}, err
	}

	err = s.store.Save(ctx, *tag)
	if err != nil {
		if err.Error() == bareknews.ErrDataAlreadyExist.Error() {
			return Response{}, bareknews.ErrDataAlreadyExist
		}
		return Response{}, err
	}

	return Response{
		ID:   tag.Label.ID,
		Name: tag.Label.Name,
		Slug: tag.Slug.String(),
	}, nil
}

func (s Service) Update(ctx context.Context, id uuid.UUID, newTagname string) (Response, error) {
	tag, err := s.store.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Response{}, bareknews.ErrDataNotFound
		} else {
			return Response{}, err
		}
	}

	tag.ChangeName(strings.TrimSpace(newTagname))

	err = tag.Validate()
	if err != nil {
		return Response{}, err
	}

	err = s.store.Update(ctx, *tag)
	if err != nil {
		return Response{}, err
	}

	return Response{
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

func (s Service) GetById(ctx context.Context, id uuid.UUID) (Response, error) {
	tg, err := s.store.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Response{}, bareknews.ErrDataNotFound
		} else {
			return Response{}, err
		}
	}

	return Response{
		ID:   tg.Label.ID,
		Name: tg.Label.Name,
		Slug: tg.Slug.String(),
	}, nil
}

func (s Service) GetByIds(ctx context.Context, ids []uuid.UUID) ([]Response, error) {
	tgs, err := s.store.GetByIds(ctx, ids)
	if err != nil {
		return []Response{}, bareknews.ErrInternalServer
	}

	r := make([]Response, 0)

	for _, t := range tgs {
		r = append(r, Response{
			ID:   t.Label.ID,
			Name: t.Label.Name,
			Slug: t.Slug.String(),
		})
	}

	return r, nil
}

func (s Service) GetAll(ctx context.Context) ([]Response, error) {
	tg, err := s.store.GetAll(ctx)
	if err != nil {
		return []Response{}, err
	}

	r := make([]Response, 0)

	for _, t := range tg {
		r = append(r, Response{
			ID:   t.Label.ID,
			Name: t.Label.Name,
			Slug: t.Slug.String(),
		})
	}

	return r, nil
}

func (s Service) GetByNames(ctx context.Context, names []string) []Response {
	tg, err := s.store.GetByNames(ctx, names...)
	if err != nil {
		return []Response{}
	}

	r := make([]Response, 0)

	for _, t := range tg {
		r = append(r, Response{
			ID:   t.Label.ID,
			Name: t.Label.Name,
			Slug: t.Slug.String(),
		})
	}

	return r
}

func (s Service) GetByName(ctx context.Context, name string) Response {
	tg, err := s.store.GetByName(ctx, name)
	if err != nil {
		return Response{}
	}

	return Response{
		ID:   tg.Label.ID,
		Name: tg.Label.Name,
		Slug: tg.Slug.String(),
	}
}