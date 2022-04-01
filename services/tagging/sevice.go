package tagging

import (
	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/domain/tag"
)

type Service struct {
	storage tag.Repository
}

func New(repo tag.Repository) Service {
	return Service{repo}
}

func (s Service) Create(tagName string) (bareknews.Tags, error) {
	tag, err := bareknews.NewTags(tagName)
	if err != nil {
		return bareknews.Tags{}, err
	}

	err = tag.Validate()
	if err != nil {
		return bareknews.Tags{}, err
	}

	err = s.storage.Save(*tag)
	if err != nil {
		return bareknews.Tags{}, err
	}

	t, err := s.storage.GetByName(tagName)
	if err != nil {
		return bareknews.Tags{}, err
	}

	return t, nil
}

func (s Service) Update(id, newTagname string) error {
	tag, err := s.storage.GetById(id)
	if err != nil {
		return err
	}

	tag.ChangeName(newTagname)
	
	err = tag.Validate()
	if err != nil {
		return err
	}

	err = s.storage.Update(*tag)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Delete(id string) error {
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

func (s Service) GetById(id string) (bareknews.Tags, error) {
	tagsResult, err := s.storage.GetById(id)
	if err != nil {
		return bareknews.Tags{}, err
	}

	return *tagsResult, nil
}

func (s Service) GetAll() ([]bareknews.Tags, error) {
	tagsResults, err := s.storage.GetAll()
	if err != nil {
		return []bareknews.Tags{}, err
	}

	return tagsResults, nil
}

// func (s Service) GetByNewsId(newsId string) ([]bareknews.Tags, error) {
// 	result, err := s.storage.GetByNewsId(newsId)
// 	if err != nil {
// 		return []bareknews.Tags{}, err
// 	}

// 	return result, nil
// }