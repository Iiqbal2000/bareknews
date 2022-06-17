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

type NewsIn struct {
	Title  string   `json:"title" validate:"required"`
	Body   string   `json:"body" validate:"required"`
	Status string   `json:"status" enums:"publish,draft" default:"draft"`
	Tags   []string `json:"tags"`
}

type NewsOut struct {
	ID     uuid.UUID       `json:"id"`
	Title  string          `json:"title"`
	Body   string          `json:"body"`
	Status string          `json:"status"`
	Slug   string          `json:"slug"`
	Tags   []tags.Response `json:"tags"`
}

func createNewsOut(n *News, tgs []tags.Response) NewsOut {
	return NewsOut{
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

func (s Service) Create(ctx context.Context, input NewsIn) (NewsOut, error) {
	tg := s.taggingSvc.GetByNames(ctx, input.Tags)
	tgId := make([]uuid.UUID, 0)

	for _, t := range tg {
		tgId = append(tgId, t.ID)
	}

	news := Create(input.Title, input.Body, bareknews.Status(input.Status), tgId)
	err := news.Validate()
	if err != nil {
		return NewsOut{}, err
	}

	err = s.storage.Save(ctx, *news)
	if err != nil {
		return NewsOut{}, err
	}

	return createNewsOut(news, tg), nil
}

func (s Service) Update(ctx context.Context, id uuid.UUID, input NewsIn) (NewsOut, error) {
	news, err := s.storage.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewsOut{}, bareknews.ErrDataNotFound
		} else {
			return NewsOut{}, err
		}
	}

	if input.Title != "" && strings.TrimSpace(input.Title) != "" {
		news.ChangeTitle(input.Title)
	}

	if input.Body != "" && strings.TrimSpace(input.Body) != "" {
		news.ChangeBody(input.Body)
	}

	if input.Status != "" && strings.TrimSpace(input.Status) != "" {
		news.ChangeStatus(bareknews.Status(input.Status))
	}

	if len(input.Tags) > 0 {
		tg := s.taggingSvc.GetByNames(ctx, input.Tags)
		tgId := make([]uuid.UUID, 0)

		for _, t := range tg {
			tgId = append(tgId, t.ID)
		}

		news.ChangeTags(tgId)
	}

	err = news.Validate()
	if err != nil {
		return NewsOut{}, err
	}

	err = s.storage.Update(ctx, *news)
	if err != nil {
		return NewsOut{}, err
	}

	tg, err := s.taggingSvc.GetByIds(ctx, news.TagsID)
	if err != nil {
		return NewsOut{}, err
	}

	return createNewsOut(news, tg), nil
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

func (s Service) GetById(ctx context.Context, id uuid.UUID) (NewsOut, error) {
	news, err := s.storage.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewsOut{}, bareknews.ErrDataNotFound
		} else {
			return NewsOut{}, err
		}
	}

	tgs, err := s.taggingSvc.GetByIds(ctx, news.TagsID)
	if err != nil {
		return NewsOut{}, err
	}

	return createNewsOut(news, tgs), nil
}

func (s Service) GetAll(ctx context.Context) ([]NewsOut, error) {
	nws, err := s.storage.GetAll(ctx)
	if err != nil {
		return []NewsOut{}, err
	}

	r := make([]NewsOut, 0)

	for _, nw := range nws {
		tgs, err := s.taggingSvc.GetByIds(ctx, nw.TagsID)
		if err != nil {
			return []NewsOut{}, err
		}

		r = append(r, createNewsOut(&nw, tgs))
	}

	return r, nil
}

func (s Service) GetAllByTopic(ctx context.Context, topic string) ([]NewsOut, error) {
	tg := s.taggingSvc.GetByNames(ctx, []string{topic})

	nws, err := s.storage.GetAllByTopic(ctx, tg[0].ID)
	if err != nil {
		return []NewsOut{}, err
	}

	r := make([]NewsOut, 0)

	for _, nw := range nws {
		tgs, err := s.taggingSvc.GetByIds(ctx, nw.TagsID)
		if err != nil {
			return []NewsOut{}, err
		}

		r = append(r, createNewsOut(&nw, tgs))
	}

	return r, nil
}

func (s Service) GetAllByStatus(ctx context.Context, status string) ([]NewsOut, error) {
	stat := bareknews.Status(status)
	err := stat.Validate()
	if err != nil {
		return []NewsOut{}, err
	}

	nws, err := s.storage.GetAllByStatus(ctx, stat)
	if err != nil {
		return []NewsOut{}, err
	}

	r := make([]NewsOut, 0)

	for _, nw := range nws {
		tgs, err := s.taggingSvc.GetByIds(ctx, nw.TagsID)
		if err != nil {
			return []NewsOut{}, err
		}

		r = append(r, createNewsOut(&nw, tgs))
	}

	return r, nil
}
