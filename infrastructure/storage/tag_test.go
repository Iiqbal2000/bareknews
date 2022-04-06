package storage

import (
	"database/sql"
	"testing"

	"github.com/Iiqbal2000/bareknews/domain/tags"
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

// func TestGetAllByNames(t *testing.T) {
// 	conn := Run(dbfile, true)
// 	storage := Tag{conn}
// 	result, err := storage.GetByNames("tag 1", "tag 23")
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	fmt.Println(result)
// }

// func TestGetById(t *testing.T) {
// 	conn := Run(dbfile, true)
// 	storage := Tag{conn}
// 	_, err := storage.GetById("50b12s9b6-0497-4b75-8bf9-26b05f4e74c4")
// 	if err == nil {
// 		t.Fatal(err.Error())
// 	}
// 	fmt.Println(err)
// }
