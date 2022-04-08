package posting

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/Iiqbal2000/bareknews/services"
	"github.com/Iiqbal2000/bareknews/services/tagging"
	"github.com/google/uuid"
)

type Response struct {
	ID     uuid.UUID          `json:"id"`
	Title  string             `json:"title"`
	Body   string             `json:"body"`
	Status string             `json:"status"`
	Slug   string             `json:"slug"`
	Tags   []tagging.Response `json:"tags"`
}

type Service struct {
	storage    news.Repository
	taggingSvc tagging.Service
}

func New(repo news.Repository, taggingSvc tagging.Service) Service {
	return Service{repo, taggingSvc}
}

func (s Service) Create(title, body, status string, tagsIn []string) error {
	status = strings.ToLower(status)
	tg := s.taggingSvc.GetByNames(tagsIn)
	tgId := make([]uuid.UUID, 0)

	for _, t := range tg {
		tgId = append(tgId, t.ID)
	}

	news := news.New(title, body, domain.Status(status), tgId)
	err := news.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Save(*news)
	if err != nil {
		return services.ErrInternalServer
	}

	return nil
}

func (s Service) Update(id uuid.UUID, title, body, status string, tgIn []string) error {
	news, err := s.storage.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		} else {
			return services.ErrInternalServer
		}
	}

	if title != "" && strings.TrimSpace(title) != "" {
		news.ChangeTitle(title)
	}

	if body != "" && strings.TrimSpace(body) != "" {
		news.ChangeBody(body)
	}

	if status != "" && strings.TrimSpace(status) != "" {
		news.ChangeStatus(domain.Status(status))
	}

	if len(tgIn) > 0 {
		tg := s.taggingSvc.GetByNames(tgIn)
		tgId := make([]uuid.UUID, 0)

		for _, t := range tg {
			tgId = append(tgId, t.ID)
		}
		news.ChangeTags(tgId)
	}

	err = news.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Update(*news)
	if err != nil {
		return services.ErrInternalServer
	}

	return nil
}

func (s Service) Delete(id uuid.UUID) error {
	_, err := s.storage.Count(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		} else {
			return services.ErrInternalServer
		}
	}

	err = s.storage.Delete(id)
	if err != nil {
		return services.ErrInternalServer
	}

	return nil
}

func (s Service) GetById(id uuid.UUID) (Response, error) {
	result, err := s.storage.GetById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Response{}, err
		} else {
			return Response{}, services.ErrInternalServer
		}
	}

	tgs, err := s.taggingSvc.GetByIds(result.TagsID)
	if err != nil {
		return Response{}, err
	}

	r := Response{
		ID:     result.Post.ID,
		Title:  result.Post.Title,
		Body:   result.Post.Body,
		Status: result.Status.String(),
		Slug:   result.Slug.String(),
		Tags:   tgs,
	}

	return r, nil
}

func (s Service) GetAll() ([]Response, error) {
	nws, err := s.storage.GetAll()
	if err != nil {
		return []Response{}, services.ErrInternalServer
	}

	r := make([]Response, 0)
	
	for _, nw := range nws {
		tgs, err := s.taggingSvc.GetByIds(nw.TagsID)
		if err != nil {
			return []Response{}, err
		}

		r = append(r, Response{
			ID:     nw.Post.ID,
			Title:  nw.Post.Title,
			Body:   nw.Post.Body,
			Status: nw.Status.String(),
			Slug:   nw.Slug.String(),
			Tags:   tgs,
		})
	}

	return r, nil
}
