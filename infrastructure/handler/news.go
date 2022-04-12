package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/services/posting"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type News struct {
	Service posting.Service
}

type NewsInput struct {
	Title  string   `json:"title"`
	Status string   `json:"status"`
	Tags   []string `json:"tags"`
	Body   string   `json:"body"`
}

func (n News) Route(r chi.Router) {
	r.Post("/", n.Create)
	r.Get("/{newsId}", n.GetById)
	r.Put("/{newsId}", n.Update)
	r.Delete("/{newsId}", n.Delete)
	r.Get("/", n.GetAll)
}

func (n News) Create(w http.ResponseWriter, r *http.Request) {
	payloadIn := NewsInput{}

	err := json.NewDecoder(r.Body).Decode(&payloadIn)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInvalidJSON)
		if err != nil {
			log.Println("(error) news.handler.create: ", err.Error())
		}
		return
	}

	nws, err := n.Service.Create(payloadIn.Title, payloadIn.Body, payloadIn.Status, payloadIn.Tags)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) news.handler.create: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully creating a news",
		"data":    nws,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.create: ", err.Error())
	}
}

func (n News) GetById(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) news.handler.getById: ", err.Error())
		}
		return
	}

	nws, err := n.Service.GetById(id)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) news.handler.getById: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully getting a news",
		"data":    nws,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.getById: ", err.Error())
	}
}

func (n News) Update(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) news.handler.update: ", err.Error())
		}
		return
	}

	payloadIn := NewsInput{}

	err = json.NewDecoder(r.Body).Decode(&payloadIn)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrInvalidJSON)
		if err != nil {
			log.Println("(error) news.handler.update: ", err.Error())
		}
		return
	}

	nws, err := n.Service.Update(
		id,
		payloadIn.Title,
		payloadIn.Body,
		payloadIn.Status,
		payloadIn.Tags,
	)

	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) news.handler.update: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully updating a news",
		"data":    nws,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.update: ", err.Error())
	}
}

func (n News) Delete(w http.ResponseWriter, r *http.Request) {
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = bareknews.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) news.handler.delete: ", err.Error())
		}
		return
	}

	err = n.Service.Delete(id)
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) news.handler.delete: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully deleting a news",
		"data":    struct{}{},
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.delete: ", err.Error())
	}
}

func (n News) GetAll(w http.ResponseWriter, r *http.Request) {
	newsRes, err := n.Service.GetAll()
	if err != nil {
		err = bareknews.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) news.handler.getAll: ", err.Error())
		}
		return
	}

	payloadRes := map[string]interface{}{
		"message": "Successfuly getting all news",
		"data": newsRes,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.getAll: ", err.Error())
	}
}