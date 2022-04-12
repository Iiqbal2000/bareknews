package storage

import (
	"database/sql"
	"testing"

	"github.com/Iiqbal2000/bareknews/domain/tags"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

func TestSave(t *testing.T) {
	conn := Run(dbfile, true)
	storage := Tag{conn}
	is := is.New(t)

	tag := tags.New("tag 1")
	err := storage.Save(*tag)
	is.NoErr(err)

	got, err := storage.GetById(tag.Label.ID)
	is.NoErr(err)
	is.True(got.Label.ID.String() != "")
	is.True(got.Label.Name != "")
	is.True(got.Slug != "")
}

func TestGetAll(t *testing.T) {
	conn := Run(dbfile, true)
	storage := Tag{conn}
	is := is.New(t)
	tag1 := tags.New("tag 1")
	tag2 := tags.New("tag 2")
	err := storage.Save(*tag1)
	is.NoErr(err)
	err = storage.Save(*tag2)
	is.NoErr(err)

	got, err := storage.GetAll()
	is.NoErr(err)
	is.Equal(got[0].Label.ID, tag1.Label.ID)
	is.Equal(got[0].Label.Name, tag1.Label.Name)
	is.Equal(got[0].Slug, tag1.Slug)

	is.Equal(got[1].Label.ID, tag2.Label.ID)
	is.Equal(got[1].Label.Name, tag2.Label.Name)
	is.Equal(got[1].Slug, tag2.Slug)
}

func TestUpdate(t *testing.T) {
	conn := Run(dbfile, true)
	storage := Tag{conn}
	is := is.New(t)
	tag := tags.New("tag 1")
	err := storage.Save(*tag)
	is.NoErr(err)

	tag.ChangeName("tag 16")
	err = storage.Update(*tag)
	is.NoErr(err)
	is.Equal(tag.Label.Name, "tag 16")
	is.Equal(tag.Slug.String(), "tag-16")
}
func TestDelete(t *testing.T) {
	conn := Run(dbfile, true)
	storage := Tag{conn}
	is := is.New(t)

	tag := tags.New("tag 1")
	err := storage.Save(*tag)
	is.NoErr(err)
	err = storage.Delete(tag.Label.ID)
	is.NoErr(err)
	_, err = storage.GetById(tag.Label.ID)
	is.Equal(err, sql.ErrNoRows)
}

func TestCount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		conn := Run(dbfile, true)
		storage := Tag{conn}
		is := is.New(t)
		tag := tags.New("tag 1")
		err := storage.Save(*tag)
		is.NoErr(err)
	
		c, err := storage.Count(tag.Label.ID)
		is.NoErr(err)
		is.True(c != 0)
	})

	t.Run("error", func(t *testing.T) {
		conn := Run(dbfile, true)
		storage := Tag{conn}
		is := is.New(t)
	
		c, err := storage.Count(uuid.New())
		is.Equal(err, sql.ErrNoRows)
		is.Equal(c, 0)
	})
}

func TestGetByNames(t *testing.T) {
	conn := Run(dbfile, true)
	storage := Tag{conn}
	is := is.New(t)
	tag1 := tags.New("tag 1")
	tag2 := tags.New("tag 2")
	err := storage.Save(*tag1)
	is.NoErr(err)
	err = storage.Save(*tag2)
	is.NoErr(err)

	got, err := storage.GetByNames(tag1.Label.Name)
	is.NoErr(err)
	is.Equal(got[0].Label.ID, tag1.Label.ID)
	is.Equal(got[0].Label.Name, tag1.Label.Name)
	is.Equal(got[0].Slug, tag1.Slug)

	is.Equal(got[1].Label.ID, tag2.Label.ID)
	is.Equal(got[1].Label.Name, tag2.Label.Name)
	is.Equal(got[1].Slug, tag2.Slug)
}

func TestGetById(t *testing.T) {
	conn := Run(dbfile, true)
	storage := Tag{conn}
	is := is.New(t)
	tag1 := tags.New("tag 1")
	err := storage.Save(*tag1)
	is.NoErr(err)

	got, err := storage.GetById(tag1.Label.ID)
	is.NoErr(err)
	is.Equal(got.Label.Name, "tag 1")
}

func TestGetByIds(t *testing.T) {
	conn := Run(dbfile, true)
	storage := Tag{conn}
	is := is.New(t)
	tag1 := tags.New("tag 1")
	tag2 := tags.New("tag 2")
	err := storage.Save(*tag1)
	is.NoErr(err)
	err = storage.Save(*tag2)
	is.NoErr(err)

	got, err := storage.GetByIds([]uuid.UUID{tag1.Label.ID, tag2.Label.ID})
	is.NoErr(err)
	is.Equal(len(got), 2)
}