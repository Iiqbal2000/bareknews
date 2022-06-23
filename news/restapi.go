package news

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/pkg/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type handler struct {
	service Service
	log     *zap.SugaredLogger
}

func CreateHandler(svc Service, log *zap.SugaredLogger) handler {
	return handler{service: svc, log: log}
}

// CreateNews godoc
// @Summary      Create a news
// @Description  Create a news and return it
// @Tags         news
// @Accept       json
// @Produce      json
// @Param news body NewsIn true "A payload of new news"
// @Success      201  {object}  web.RespBody{data=posting.Response} "Response body for a new news"
// @Failure      400  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      404  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /news [post]
func (n handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payloadIn := NewsIn{}

	err := json.NewDecoder(r.Body).Decode(&payloadIn)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, bareknews.ErrInvalidJSON)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	nws, err := n.service.Create(ctx, payloadIn)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, err)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadRes := web.RespBody{
		Message: "Successfully creating a news",
		Data:    nws,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		n.log.Error(errors.Wrap(err, "failed to write a response"))
	}
}

// GetNewsById godoc
// @Summary      Get a news
// @Description  Get a news by id
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "News ID"  Format(uuid)
// @Success      200  {object}  web.RespBody{data=posting.Response} "Response body for a news"
// @Failure      404  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /news/{id} [get]
func (n handler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, bareknews.ErrDataNotFound)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	nws, err := n.service.GetById(ctx, id)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, err)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadRes := web.RespBody{
		Message: "Successfully getting a news",
		Data:    nws,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		n.log.Error(errors.Wrap(err, "failed to write a response"))
	}
}

// UpdateNews godoc
// @Summary      Update a news
// @Description  Update a news and return it
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "News ID"  Format(uuid)
// @Param news body NewsIn true "A payload of new news"
// @Success      200  {object}  web.RespBody{data=posting.Response} "Response body for a new news"
// @Failure      400  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      404  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /news/{id} [put]
func (n handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, bareknews.ErrDataNotFound)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadIn := NewsIn{}

	err = json.NewDecoder(r.Body).Decode(&payloadIn)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, bareknews.ErrInvalidJSON)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	nws, err := n.service.Update(ctx, id, payloadIn)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, err)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadRes := web.RespBody{
		Message: "Successfully updating a news",
		Data:    nws,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		n.log.Error(errors.Wrap(err, "failed to write a response"))
	}
}

// DeleteNews godoc
// @Summary      Delete a news
// @Description  Delete a news by id
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "News ID"  Format(uuid)
// @Success      200  {object}  web.RespBody{data=object}
// @Failure      404  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /news/{id} [delete]
func (n handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, bareknews.ErrDataNotFound)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	err = n.service.Delete(ctx, id)
	if err != nil {
		err = web.WriteErrResponse(w, n.log, err)
		if err != nil {
			n.log.Error(errors.Wrap(err, "failed to write a response"))
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
		n.log.Error(errors.Wrap(err, "failed to write a response"))
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
// @Success      200  {object}  web.RespBody{data=[]posting.Response} "Array of news body"
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /news [get]
func (n handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := r.URL.Query()
	topic := strings.TrimSpace(q.Get("topic"))
	status := strings.TrimSpace(q.Get("status"))

	newsRes := make([]NewsOut, 0)

	switch {
	case topic != "" && status != "":
		nws, err := n.service.GetAllByTopic(ctx, topic)
		if err != nil {
			err = web.WriteErrResponse(w, n.log, err)
			if err != nil {
				n.log.Error(errors.Wrap(err, "failed to write a response"))
			}
			return
		}

		for _, n := range nws {
			if n.Status == status {
				newsRes = append(newsRes, n)
			}
		}
	case topic == "" && status != "":
		nws, err := n.service.GetAllByStatus(ctx, status)
		if err != nil {
			err = web.WriteErrResponse(w, n.log, err)
			if err != nil {
				n.log.Error(errors.Wrap(err, "failed to write a response"))
			}
			return
		}

		newsRes = append(newsRes, nws...)
	case topic != "" && status == "":
		nws, err := n.service.GetAllByTopic(ctx, topic)
		if err != nil {
			err = web.WriteErrResponse(w, n.log, err)
			if err != nil {
				n.log.Error(errors.Wrap(err, "failed to write a response"))
			}
			return
		}

		newsRes = append(newsRes, nws...)
	default:
		nws, err := n.service.GetAll(ctx)
		if err != nil {
			err = web.WriteErrResponse(w, n.log, err)
			if err != nil {
				n.log.Error(errors.Wrap(err, "failed to write a response"))
			}
			return
		}

		newsRes = append(newsRes, nws...)
	}

	payloadRes := web.RespBody{
		Message: "Successfuly getting all news",
		Data:    newsRes,
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		n.log.Error(errors.Wrap(err, "failed to write a response"))
	}
}
