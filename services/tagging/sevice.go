package tagging

import "github.com/Iiqbal2000/bareknews/domain/tags"

type Service struct {
	storage tags.Repository
}

// func New(repo tag.Repository) Service {
// 	return Service{repo}
// }

// func (s Service) Create(tagName string) (domain.Tags, error) {
// 	tag, err := domain.NewTags(tagName)
// 	if err != nil {
// 		return domain.Tags{}, err
// 	}

// 	err = tag.Validate()
// 	if err != nil {
// 		return domain.Tags{}, err
// 	}

// 	err = s.storage.Save(*tag)
// 	if err != nil {
// 		return domain.Tags{}, err
// 	}

// 	t, err := s.storage.GetByName(tagName)
// 	if err != nil {
// 		return domain.Tags{}, err
// 	}

// 	return t, nil
// }

// func (s Service) Update(id, newTagname string) error {
// 	tag, err := s.storage.GetById(id)
// 	if err != nil {
// 		return err
// 	}

// 	tag.ChangeName(newTagname)

// 	err = tag.Validate()
// 	if err != nil {
// 		return err
// 	}

// 	err = s.storage.Update(*tag)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Service) Delete(id string) error {
// 	_, err := s.storage.GetById(id)
// 	if err != nil {
// 		return err
// 	}

// 	err = s.storage.Delete(id)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s Service) GetById(id string) (domain.Tags, error) {
// 	tagsResult, err := s.storage.GetById(id)
// 	if err != nil {
// 		return domain.Tags{}, err
// 	}

// 	return *tagsResult, nil
// }

// func (s Service) GetAll() ([]domain.Tags, error) {
// 	tagsResults, err := s.storage.GetAll()
// 	if err != nil {
// 		return []domain.Tags{}, err
// 	}

// 	return tagsResults, nil
// }

// func (s Service) GetByNames(tagsIn []string) []domain.Tags {
// 	tags, err := s.storage.GetByNames(tagsIn...)
// 	if err != nil {
// 		log.Println("Error in tagging service: ", err.Error())
// 	}
// 	return tags
// }
