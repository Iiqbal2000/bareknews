package posting

import (
	"log"
	"strings"

	"github.com/Iiqbal2000/bareknews/domain"
	"github.com/Iiqbal2000/bareknews/domain/news"
	"github.com/Iiqbal2000/bareknews/services/tagging"
	"github.com/google/uuid"
)

type Service struct {
	storage    news.Repository
	taggingSvc tagging.Service
}

func New(repo news.Repository, taggingSvc tagging.Service) Service {
	return Service{repo, taggingSvc}
}

func (s Service) Create(title, body, status string, tagsIn []string) error {
	status = strings.ToLower(status)
	tags := s.taggingSvc.GetByNames(tagsIn)
	news := news.New(title, body, domain.Status(status), tags)
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

func (s Service) GetById(id string) (news.News, error) {
	result, err := s.storage.GetById(id)
	if err != nil {
		log.Println(err.Error())
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

func (s Service) Update(id uuid.UUID, title, body, status string, tags []domain.Tags) error {
	news, err := s.storage.GetById(id.String())
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
		news.ChangeStatus(domain.Status(status))
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
	_, err := s.storage.GetById(id.String())
	if err != nil {
		return err
	}

	err = s.storage.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
