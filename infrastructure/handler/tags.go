package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/services/tagging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Tags struct {
	Service tagging.Service
}

func (t Tags) Route(r chi.Router) {
	r.Post("/", t.Create)
	r.Get("/{tagId}", t.GetById)
	r.Put("/{tagId}", t.Update)
	r.Get("/", t.GetAll)
	r.Delete("/{tagId}", t.Delete)
}

func (t Tags) Create(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Name string `json:"name"`
	}

	payload := Input{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInternalServer)
		if err != nil {
			log.Println("(error) tags.handler.create: ", err.Error())
		}
		return
	}

	tagRes, err := t.Service.Create(payload.Name)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.create: ", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Successfully creating a tag",
		"data": tagRes,
	})
	if err != nil {
		log.Println("(error) tags.handler.create: ", err.Error())
	}
}

func (t Tags) GetById(w http.ResponseWriter, r *http.Request) {
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInternalServer)
		if err != nil {
			log.Println("(error) tags.handler.getById: ", err.Error())
		}
		return
	}

	tg, err := t.Service.GetById(id)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.getById: ", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Successfully getting a tag",
		"data": tg,
	})
	if err != nil {
		log.Println("(error) tags.handler.getById: ", err.Error())
	}
}

func (t Tags) Update(w http.ResponseWriter, r *http.Request) {
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInternalServer)
		if err != nil {
			log.Println("(error) tags.handler.update: ", err.Error())
		}
		return
	}

	type Input struct {
		Name string `json:"name"`
	}

	payload := Input{}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInternalServer)
		if err != nil {
			log.Println("(error) tags.handler.update: ", err.Error())
		}
		return
	}

	tg, err := t.Service.Update(id, payload.Name)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.update: ", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Successfully updating a tag",
		"data": tg,
	})
	if err != nil {
		log.Println("(error) tags.handler.update: ", err.Error())
	}
}

func (t Tags) GetAll(w http.ResponseWriter, r *http.Request) {
	tgs, err := t.Service.GetAll()
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInternalServer)
		if err != nil {
			log.Println("(error) tags.handler.getAll: ", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Successfully getting all tags",
		"data": tgs,
	})
	if err != nil {
		log.Println("(error) tags.handler.getAll: ", err.Error())
	}
}

func (t Tags) Delete(w http.ResponseWriter, r *http.Request) {
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInternalServer)
		if err != nil {
			log.Println("(error) tags.handler.delete: ", err.Error())
		}
		return
	}

	err = t.Service.Delete(id)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.delete: ", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Successfully deleting a tag",
		"data": struct{}{},
	})
	if err != nil {
		log.Println("(error) tags.handler.getById: ", err.Error())
	}
}