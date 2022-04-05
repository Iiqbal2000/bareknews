package tags

//go:generate moq -out tagRepo_moq.go . Repository
type Repository interface {
	Save(Tags) error
	Update(Tags) error
	GetById(string) (*Tags, error)
	GetByNewsId(id string) ([]Tags, error)
	GetByName(names string) (Tags, error)
	GetByNames(names ...string) ([]Tags, error)
	GetAll() ([]Tags, error)
	Delete(id string) error
}
