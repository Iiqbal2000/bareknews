package storage

import (
	"fmt"
	"testing"

	"github.com/Iiqbal2000/bareknews"
)

func TestSave(t *testing.T) {
	conn := Run()
	storage := TagStorage{conn}
	tag, err := bareknews.NewTags("tag 1")
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Save(*tag)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	conn := Run()
	storage := TagStorage{conn}
	tag, err := storage.GetById("f83907d-ac88-45de-819c-809445b03d14")
	if err != nil {
		t.Fatal(err)
	}
	tag.ChangeName("tag 16")
	err = storage.Update(*tag)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAll(t *testing.T) {
	conn := Run()
	storage := TagStorage{conn}
	tag1, _ := bareknews.NewTags("tag 1")
	tag2, _ := bareknews.NewTags("tag 2")
	err := storage.Save(*tag1)
	if err != nil {
		t.Fatal(err)
	}
	err = storage.Save(*tag2)
	if err != nil {
		t.Fatal(err)
	}
	result, err := storage.GetAll()
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(result)
}

func TestGetAllByNames(t *testing.T) {
	conn := Run()
	storage := TagStorage{conn}
	result, err := storage.GetByNames("tag 1", "tag 23")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(result)
}

func TestGetById(t *testing.T) {
	conn := Run()
	storage := TagStorage{conn}
	_, err := storage.GetById("50b12s9b6-0497-4b75-8bf9-26b05f4e74c4")
	if err == nil {
		t.Fatal(err.Error())
	}
	fmt.Println(err)
}

func TestDelete(t *testing.T) {
	conn := Run()
	storage := TagStorage{conn}
	err := storage.Delete("50b12s9b6-0497-4b75-8bf9-26b05f4e74c4")
	if err != nil {
		t.Fatal(err.Error())
	}

	result, err := storage.GetAll()
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(result)
}
