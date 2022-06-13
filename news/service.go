package news

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/google/uuid"
)

type Response struct {
	ID     uuid.UUID       `json:"id"`
	Title  string          `json:"title"`
	Body   string          `json:"body"`
	Status string          `json:"status"`
	Slug   string          `json:"slug"`
	Tags   []tags.Response `json:"tags"`
}

func createResponse(n *News, tgs []tags.Response) Response {
	return Response{
		ID:     n.Post.ID,
		Title:  n.Post.Title,
		Body:   n.Post.Body,
		Status: n.Status.String(),
		Slug:   n.Slug.String(),
		Tags:   tgs,
	}
}

type Service struct {
	storage    Repository
	taggingSvc tags.Service
}

func CreateSvc(repo Repository, taggingSvc tags.Service) Service {
	return Service{repo, taggingSvc}
}

func (s Service) Create(ctx context.Context, title, body, status string, tagsIn []string) (Response, error) {
	status = strings.ToLower(status)
	tg := s.taggingSvc.GetByNames(ctx, tagsIn)
	tgId := make([]uuid.UUID, 0)

	for _, t := range tg {
		tgId = append(tgId, t.ID)
	}

	news := Create(title, body, bareknews.Status(status), tgId)
	err := news.Validate()
	if err != nil {
		return Response{}, err
	}

	err = s.storage.Save(ctx, *news)
	if err != nil {
		return Response{}, err
	}

	return createResponse(news, tg), nil
}

func (s Service) Update(ctx context.Context, id uuid.UUID, title, body, status string, tgIn []string) (Response, error) {
	news, err := s.storage.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Response{}, bareknews.ErrDataNotFound
		} else {
			return Response{}, err
		}
	}

	if title != "" && strings.TrimSpace(title) != "" {
		news.ChangeTitle(title)
	}

	if body != "" && strings.TrimSpace(body) != "" {
		news.ChangeBody(body)
	}

	if status != "" && strings.TrimSpace(status) != "" {
		news.ChangeStatus(bareknews.Status(status))
	}

	if len(tgIn) > 0 {
		tg := s.taggingSvc.GetByNames(ctx, tgIn)
		tgId := make([]uuid.UUID, 0)

		for _, t := range tg {
			tgId = append(tgId, t.ID)
		}

		news.ChangeTags(tgId)
	}

	err = news.Validate()
	if err != nil {
		return Response{}, err
	}

	err = s.storage.Update(ctx, *news)
	if err != nil {
		return Response{}, err
	}

	tg, err := s.taggingSvc.GetByIds(ctx, news.TagsID)
	if err != nil {
		return Response{}, err
	}

	return createResponse(news, tg), nil
}

func (s Service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.storage.Count(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return bareknews.ErrDataNotFound
		} else {
			return err
		}
	}

	err = s.storage.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetById(ctx context.Context, id uuid.UUID) (Response, error) {
	news, err := s.storage.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Response{}, bareknews.ErrDataNotFound
		} else {
			return Response{}, err
		}
	}

	tgs, err := s.taggingSvc.GetByIds(ctx, news.TagsID)
	if err != nil {
		return Response{}, err
	}

	return createResponse(news, tgs), nil
}

func (s Service) GetAll(ctx context.Context) ([]Response, error) {
	nws, err := s.storage.GetAll(ctx)
	if err != nil {
		return []Response{}, err
	}

	r := make([]Response, 0)

	for _, nw := range nws {
		tgs, err := s.taggingSvc.GetByIds(ctx, nw.TagsID)
		if err != nil {
			return []Response{}, err
		}

		r = append(r, createResponse(&nw, tgs))
	}

	return r, nil
}

func (s Service) GetAllByTopic(ctx context.Context, topic string) ([]Response, error) {
	tg := s.taggingSvc.GetByNames(ctx, []string{topic})

	nws, err := s.storage.GetAllByTopic(ctx, tg[0].ID)
	if err != nil {
		return []Response{}, err
	}

	r := make([]Response, 0)

	for _, nw := range nws {
		tgs, err := s.taggingSvc.GetByIds(ctx, nw.TagsID)
		if err != nil {
			return []Response{}, err
		}

		r = append(r, createResponse(&nw, tgs))
	}

	return r, nil
}

func (s Service) GetAllByStatus(ctx context.Context, status string) ([]Response, error) {
	stat := bareknews.Status(status)
	err := stat.Validate()
	if err != nil {
		return []Response{}, err
	}

	nws, err := s.storage.GetAllByStatus(ctx, stat)
	if err != nil {
		return []Response{}, err
	}

	r := make([]Response, 0)

	for _, nw := range nws {
		tgs, err := s.taggingSvc.GetByIds(ctx, nw.TagsID)
		if err != nil {
			return []Response{}, err
		}

		r = append(r, createResponse(&nw, tgs))
	}

	return r, nil
}
