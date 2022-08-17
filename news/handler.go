package news

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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
func (n handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	payloadIn := NewsIn{}

	err := json.NewDecoder(r.Body).Decode(&payloadIn)
	if err != nil {
		return bareknews.ErrInvalidJSON
	}

	nws, err := n.service.Create(ctx, payloadIn)
	if err != nil {
		if errors.Is(err, bareknews.ErrDataAlreadyExist) {
			return web.NewRequestError(bareknews.ErrDataAlreadyExist, http.StatusConflict)
		}
		return err
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfully creating a news",
		Data:    nws,
	}

	return web.Respond(w, payloadRes, http.StatusCreated)
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
func (n handler) GetById(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rawID := chi.URLParam(r, "newsId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
	}

	nws, err := n.service.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
		}
		return err
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfully getting a news",
		Data:    nws,
	}

	return web.Respond(w, payloadRes, http.StatusOK)
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
func (n handler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rawID := chi.URLParam(r, "newsId")
	id, err := uuid.Parse(rawID)
	if err != nil {
		return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
	}

	payloadIn := NewsIn{}

	err = json.NewDecoder(r.Body).Decode(&payloadIn)
	if err != nil {
		return bareknews.ErrInvalidJSON
	}

	nws, err := n.service.Update(ctx, id, payloadIn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
		}
		return err
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfully updating a news",
		Data:    nws,
	}

	return web.Respond(w, payloadRes, http.StatusOK)
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
func (n handler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rawID := chi.URLParam(r, "newsId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
	}

	err = n.service.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
		}
		return err
	}

	payloadRes := map[string]interface{}{
		"message": "Successfully deleting a news",
		"data":    struct{}{},
	}

	return web.Respond(w, payloadRes, http.StatusOK)
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
func (n handler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()
	topic := strings.TrimSpace(q.Get("topic"))
	status := strings.TrimSpace(q.Get("status"))
	rawCursor := strings.TrimSpace(q.Get("cursor"))

	if rawCursor == "" {
		rawCursor = "0"
	}

	cursor, err := strconv.ParseInt(rawCursor, 10, 64)
	if err != nil {
		return web.NewRequestError(errors.New("failed to convert the cursor"), http.StatusBadRequest)
	}

	newsRes := make([]NewsOut, 0)

	switch {
	case topic != "" && status != "":
		nws, err := n.service.GetAllByTopic(ctx, topic, cursor)
		if err != nil {
			return err
		}

		for _, n := range nws {
			if n.Status == status {
				newsRes = append(newsRes, n)
			}
		}
	case topic == "" && status != "":
		nws, err := n.service.GetAllByStatus(ctx, status, cursor)
		if err != nil {
			return err
		}

		newsRes = append(newsRes, nws...)
	case topic != "" && status == "":
		nws, err := n.service.GetAllByTopic(ctx, topic, cursor)
		if err != nil {
			return err
		}

		newsRes = append(newsRes, nws...)
	default:
		nws, err := n.service.GetAll(ctx, cursor)
		if err != nil {
			return err
		}

		newsRes = append(newsRes, nws...)
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfuly getting all news",
		Data:    newsRes,
	}

	return web.Respond(w, payloadRes, http.StatusOK)
}
