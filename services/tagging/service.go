package tagging

import (
	"database/sql"
	"errors"

	"github.com/Iiqbal2000/bareknews/domain/tags"
	"github.com/Iiqbal2000/bareknews"
	"github.com/google/uuid"
)

type Response struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type Service struct {
	storage tags.Repository
}

func New(repo tags.Repository) Service {
	return Service{repo}
}

func (s Service) Create(tagName string) (Response, error) {
	tag := tags.New(tagName)

	err := tag.Validate()
	if err != nil {
		return Response{}, err
	}

	err = s.storage.Save(*tag)
	if err != nil {
		if err.Error() == bareknews.ErrDataAlreadyExist.Error() {
			return Response{}, bareknews.ErrDataAlreadyExist
		}
		return Response{}, bareknews.ErrInternalServer
	}

	return Response{
		ID:   tag.Label.ID,
		Name: tag.Label.Name,
		Slug: tag.Slug.String(),
	}, nil
}

func (s Service) Update(id uuid.UUID, newTagname string) (Response, error) {
	tag, err := s.storage.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Response{}, bareknews.ErrDataNotFound
		} else {
			return Response{}, bareknews.ErrInternalServer
		}
	}

	tag.ChangeName(newTagname)

	err = tag.Validate()
	if err != nil {
		return Response{}, err
	}

	err = s.storage.Update(*tag)
	if err != nil {
		return Response{}, bareknews.ErrInternalServer
	}

	return Response{
		ID: tag.Label.ID,
		Name: tag.Label.Name,
		Slug: tag.Slug.String(),
	}, nil
}

func (s Service) Delete(id uuid.UUID) error {
	_, err := s.storage.Count(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return bareknews.ErrDataNotFound
		} else {
			return bareknews.ErrInternalServer
		}
	}

	err = s.storage.Delete(id)
	if err != nil {
		return bareknews.ErrInternalServer
	}

	return nil
}

func (s Service) GetById(id uuid.UUID) (Response, error) {
	tg, err := s.storage.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Response{}, bareknews.ErrDataNotFound
		} else {
			return Response{}, bareknews.ErrInternalServer
		}
	}

	return Response{
		ID:   tg.Label.ID,
		Name: tg.Label.Name,
		Slug: tg.Slug.String(),
	}, nil
}

func (s Service) GetByIds(ids []uuid.UUID) ([]Response, error) {
	tgs, err := s.storage.GetByIds(ids)
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

func (s Service) GetAll() ([]Response, error) {
	tg, err := s.storage.GetAll()
	if err != nil {
		return []Response{}, bareknews.ErrInternalServer
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

func (s Service) GetByNames(tagsIn []string) []Response {
	tg, err := s.storage.GetByNames(tagsIn...)
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
