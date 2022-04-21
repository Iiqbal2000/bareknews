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
	ctx := r.Context()

	type Input struct {
		Name string `json:"name"`
	}

	payload := Input{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInvalidJSON)
		if err != nil {
			log.Println("(error) tags.handler.create: ", err.Error())
		}
		return
	}

	tagRes, err := t.Service.Create(ctx, payload.Name)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.create: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully creating a tag",
		"data": tagRes,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.create: ", err.Error())
	}
}

func (t Tags) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) tags.handler.getById: ", err.Error())
		}
		return
	}

	tg, err := t.Service.GetById(ctx, id)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.getById: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully getting a tag",
		"data": tg,
	}

	w.WriteHeader(http.StatusOK)
	
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.getById: ", err.Error())
	}
}

func (t Tags) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrDataNotFound)
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
		err = bareknews.WriteErrResponse(w, bareknews.ErrInvalidJSON)
		if err != nil {
			log.Println("(error) tags.handler.update: ", err.Error())
		}
		return
	}

	tg, err := t.Service.Update(ctx, id, payload.Name)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.update: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully updating a tag",
		"data": tg,
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.update: ", err.Error())
	}
}

func (t Tags) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tgs, err := t.Service.GetAll(ctx)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.getAll: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully getting all tags",
		"data": tgs,
	}

	w.WriteHeader(http.StatusOK)
	
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.getAll: ", err.Error())
	}
}

func (t Tags) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) tags.handler.delete: ", err.Error())
		}
		return
	}

	err = t.Service.Delete(ctx, id)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.delete: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully deleting a tag",
		"data": struct{}{},
	}

	w.WriteHeader(http.StatusOK)
	
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.delete: ", err.Error())
	}
}