package news

import (
	"context"
	"strings"
	"time"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("github.com/Iiqbal2000/bareknews/news")

type NewsIn struct {
	Title  string   `json:"title" validate:"required"`
	Body   string   `json:"body" validate:"required"`
	Status string   `json:"status" enums:"publish,draft" default:"draft"`
	Tags   []string `json:"tags"`
}

type NewsOut struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	Status      string         `json:"status"`
	Slug        string         `json:"slug"`
	Tags        []tags.TagsOut `json:"tags"`
	DateCreated int64          `json:"date_created"`
	DateUpdated int64          `json:"date_updated"`
}

func createNewsOut(n *News, tgs []tags.TagsOut) NewsOut {
	return NewsOut{
		ID:          n.Post.ID,
		Title:       n.Post.Title,
		Body:        n.Post.Body,
		Status:      n.Status.String(),
		Slug:        n.Slug.String(),
		Tags:        tgs,
		DateCreated: n.DateCreated,
		DateUpdated: n.DateUpdated,
	}
}

type Service struct {
	store   Repository
	tagging tags.Service
}

func CreateSvc(repo Repository, tagging tags.Service) Service {
	return Service{repo, tagging}
}

func (s Service) Create(ctx context.Context, input NewsIn) (NewsOut, error) {
	ctx, span := tracer.Start(ctx, "news.Create")
	defer span.End()

	tg := s.tagging.GetByNames(ctx, input.Tags)
	tgId := make([]uuid.UUID, 0)

	for _, t := range tg {
		tgId = append(tgId, t.ID)
	}

	news := Create(input.Title, input.Body, bareknews.Status(input.Status), tgId, time.Now().Unix())

	err := news.Validate()
	if err != nil {
		return NewsOut{}, err
	}

	err = s.store.Save(ctx, *news)
	if err != nil {
		return NewsOut{}, errors.Wrap(err, "save a news")
	}

	return createNewsOut(news, tg), nil
}

func (s Service) Update(ctx context.Context, id uuid.UUID, input NewsIn) (NewsOut, error) {
	ctx, span := tracer.Start(ctx, "news.Update")
	defer span.End()

	news, err := s.store.GetById(ctx, id)
	if err != nil {
		return NewsOut{}, err
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

	news.ChangeDateUpdated(time.Now().Unix())

	if len(input.Tags) > 0 {
		tg := s.tagging.GetByNames(ctx, input.Tags)
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

	err = s.store.Update(ctx, *news)
	if err != nil {
		return NewsOut{}, errors.Wrap(err, "update a news item")
	}

	tg, err := s.tagging.GetByIds(ctx, news.TagsID)
	if err != nil {
		return NewsOut{}, errors.Wrap(err, "get tags by ids")
	}

	return createNewsOut(news, tg), nil
}

func (s Service) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracer.Start(ctx, "news.Delete")
	defer span.End()

	_, err := s.store.Count(ctx, id)
	if err != nil {
		return err
	}

	err = s.store.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "delete a news item")
	}

	return nil
}

func (s Service) GetById(ctx context.Context, id uuid.UUID) (NewsOut, error) {
	ctx, span := tracer.Start(ctx, "news.GetById")
	defer span.End()

	news, err := s.store.GetById(ctx, id)
	if err != nil {
		return NewsOut{}, err
	}

	tgs, err := s.tagging.GetByIds(ctx, news.TagsID)
	if err != nil {
		return NewsOut{}, errors.Wrap(err, "get tags by ids")
	}

	return createNewsOut(news, tgs), nil
}

func (s Service) GetAll(ctx context.Context, cursor int64) ([]NewsOut, error) {
	ctx, span := tracer.Start(ctx, "news.GetAll")
	defer span.End()

	nws, err := s.store.GetAll(ctx, cursor, 2)
	if err != nil {
		return []NewsOut{}, errors.Wrap(err, "get all news items")
	}

	r := make([]NewsOut, 0)

	for _, nw := range nws {
		tgs, err := s.tagging.GetByIds(ctx, nw.TagsID)
		if err != nil {
			return []NewsOut{}, errors.Wrap(err, "get tags by ids")
		}

		r = append(r, createNewsOut(&nw, tgs))
	}

	return r, nil
}

func (s Service) GetAllByTopic(ctx context.Context, topic string, cursor int64) ([]NewsOut, error) {
	ctx, span := tracer.Start(ctx, "news.GetAllByTopic")
	defer span.End()

	tg := s.tagging.GetByName(ctx, topic)

	newsItems, err := s.store.GetAllByTopic(ctx, tg.ID, cursor, 2)
	if err != nil {
		return []NewsOut{}, errors.Wrap(err, "get all news items by topic")
	}

	r := make([]NewsOut, 0)

	for _, item := range newsItems {
		tgs, err := s.tagging.GetByIds(ctx, item.TagsID)
		if err != nil {
			return []NewsOut{}, errors.Wrap(err, "get tags by ids")
		}

		r = append(r, createNewsOut(&item, tgs))
	}

	return r, nil
}

func (s Service) GetAllByStatus(ctx context.Context, statusIn string, cursor int64) ([]NewsOut, error) {
	ctx, span := tracer.Start(ctx, "news.GetAllByStatus")
	defer span.End()

	status := bareknews.Status(statusIn)
	err := status.Validate()
	if err != nil {
		return []NewsOut{}, err
	}

	nws, err := s.store.GetAllByStatus(ctx, status, cursor, 2)
	if err != nil {
		return []NewsOut{}, errors.Wrap(err, "get all news items by status")
	}

	r := make([]NewsOut, 0)

	for _, nw := range nws {
		tgs, err := s.tagging.GetByIds(ctx, nw.TagsID)
		if err != nil {
			return []NewsOut{}, errors.Wrap(err, "get tags by ids")
		}

		r = append(r, createNewsOut(&nw, tgs))
	}

	return r, nil
}
