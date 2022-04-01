package posting

import (
	"strings"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/google/uuid"
)

type Service struct {
	storage news.Repository
}

func New(repo news.Repository) Service {
	return Service{repo}
}

func (s Service) Create(title, body, status string, tags []bareknews.Tags) error {
	status = strings.ToLower(status)
	news := news.New(title, body, bareknews.Status(status), tags)
	err := news.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Save(*news)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetById(id uuid.UUID) (news.News, error) {
	result, err := s.storage.GetById(id)
	if err != nil {
		return news.News{}, err
	}

	return *result, nil
}

func (s Service) GetAll() ([]news.News, error) {
	result, err := s.storage.GetAll()
	if err != nil {
		return []news.News{}, err
	}

	return result, nil
}

func (s Service) Update(id uuid.UUID, title, body, status string, tags []bareknews.Tags) error {
	news, err := s.storage.GetById(id)
	if err != nil {
		return err
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

	if len(tags) > 0 {
		news.ChangeTags(tags)
	}

	err = news.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Update(*news)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Delete(id uuid.UUID) error {
	_, err := s.storage.GetById(id)
	if err != nil {
		return err
	}

	err = s.storage.Delete(id)
	if err != nil {
		return err
	}

	return nil
}