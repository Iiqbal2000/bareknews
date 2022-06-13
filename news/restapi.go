package news

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/pkg/restapi"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Restapi struct {
	Service Service
}

type InputNews struct {
	Title  string   `json:"title" validate:"required"`
	Status string   `json:"status" enums:"publish,draft" default:"draft"`
	Tags   []string `json:"tags"`
	Body   string   `json:"body" validate:"required"`
}

func (n Restapi) Route(r chi.Router) {
	r.Post("/", n.Create)
	r.Get("/{newsId}", n.GetById)
	r.Put("/{newsId}", n.Update)
	r.Delete("/{newsId}", n.Delete)
	r.Get("/", n.GetAll)
}

// CreateNews godoc
// @Summary      Create a news
// @Description  Create a news and return it
// @Tags         news
// @Accept       json
// @Produce      json
// @Param news body InputNews true "A payload of new news"
// @Success      201  {object}  restapi.RespBody{data=posting.Response} "Response body for a new news"
// @Failure      400  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      404  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /news [post]
func (n Restapi) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payloadIn := InputNews{}

	err := json.NewDecoder(r.Body).Decode(&payloadIn)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrInvalidJSON)
		if err != nil {
			log.Println("(error) news.handler.create: ", err.Error())
		}
		return
	}

	nws, err := n.Service.Create(ctx, payloadIn.Title, payloadIn.Body, payloadIn.Status, payloadIn.Tags)
	if err != nil {
		err = restapi.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) news.handler.create: ", err.Error())
		}
		return
	}

	payloadRes := restapi.RespBody{
		Message: "Successfully creating a news",
		Data:    nws,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.create: ", err.Error())
	}
}

// GetNewsById godoc
// @Summary      Get a news
// @Description  Get a news by id
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "News ID"  Format(uuid)
// @Success      200  {object}  restapi.RespBody{data=posting.Response} "Response body for a news"
// @Failure      404  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /news/{id} [get]
func (n Restapi) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) news.handler.getById: ", err.Error())
		}
		return
	}

	nws, err := n.Service.GetById(ctx, id)
	if err != nil {
		err = restapi.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) news.handler.getById: ", err.Error())
		}
		return
	}

	payloadRes := restapi.RespBody{
		Message: "Successfully getting a news",
		Data:    nws,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.getById: ", err.Error())
	}
}

// UpdateNews godoc
// @Summary      Update a news
// @Description  Update a news and return it
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "News ID"  Format(uuid)
// @Param news body InputNews true "A payload of new news"
// @Success      200  {object}  restapi.RespBody{data=posting.Response} "Response body for a new news"
// @Failure      400  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      404  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /news/{id} [put]
func (n Restapi) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) news.handler.update: ", err.Error())
		}
		return
	}

	payloadIn := InputNews{}

	err = json.NewDecoder(r.Body).Decode(&payloadIn)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrInvalidJSON)
		if err != nil {
			log.Println("(error) news.handler.update: ", err.Error())
		}
		return
	}

	nws, err := n.Service.Update(ctx,
		id,
		payloadIn.Title,
		payloadIn.Body,
		payloadIn.Status,
		payloadIn.Tags,
	)

	if err != nil {
		err = restapi.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) news.handler.update: ", err.Error())
		}
		return
	}

	payloadRes := restapi.RespBody{
		Message: "Successfully updating a news",
		Data:    nws,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.update: ", err.Error())
	}
}

// DeleteNews godoc
// @Summary      Delete a news
// @Description  Delete a news by id
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "News ID"  Format(uuid)
// @Success      200  {object}  restapi.RespBody{data=object}
// @Failure      404  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /news/{id} [delete]
func (n Restapi) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) news.handler.delete: ", err.Error())
		}
		return
	}

	err = n.Service.Delete(ctx, id)
	if err != nil {
		err = restapi.WriteErrResponse(w, err)
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

// GetAllNews godoc
// @Summary      Get all news
// @Description  Get all news
// @Tags         news
// @Accept       json
// @Produce      json
// @Param   topic      query     string     false  "a topic"
// @Param   status      query     string     false  "status of the news"	Enums(draft, publish)
// @Success      200  {object}  restapi.RespBody{data=[]posting.Response} "Array of news body"
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /news [get]
func (n Restapi) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()
	topic := q.Get("topic")
	status := q.Get("status")

	newsRes := make([]Response, 0)

	switch {
	case topic != "" && status != "":
		nws, err := n.Service.GetAllByTopic(ctx, topic)
		if err != nil {
			err = restapi.WriteErrResponse(w, err)
			if err != nil {
				log.Println("(error) news.handler.getAll: ", err.Error())
			}
			return
		}

		for _, n := range nws {
			if n.Status == status {
				newsRes = append(newsRes, n)
			}
		}
	case topic == "" && status != "":
		nws, err := n.Service.GetAllByStatus(ctx, status)
		if err != nil {
			err = restapi.WriteErrResponse(w, err)
			if err != nil {
				log.Println("(error) news.handler.getAll: ", err.Error())
			}
			return
		}

		newsRes = append(newsRes, nws...)
	case topic != "" && status == "":
		nws, err := n.Service.GetAllByTopic(ctx, topic)
		if err != nil {
			err = restapi.WriteErrResponse(w, err)
			if err != nil {
				log.Println("(error) news.handler.getAll: ", err.Error())
			}
			return
		}

		newsRes = append(newsRes, nws...)
	default:
		nws, err := n.Service.GetAll(ctx)
		if err != nil {
			err = restapi.WriteErrResponse(w, err)
			if err != nil {
				log.Println("(error) news.handler.getAll: ", err.Error())
			}
			return
		}

		newsRes = append(newsRes, nws...)
	}

	payloadRes := restapi.RespBody{
		Message: "Successfuly getting all news",
		Data:    newsRes,
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) news.handler.getAll: ", err.Error())
	}
}
