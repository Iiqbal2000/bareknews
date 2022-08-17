package tags

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Iiqbal2000/bareknews"
	"github.com/Iiqbal2000/bareknews/pkg/web"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type handler struct {
	service Service
	log     *zap.SugaredLogger
}

type InputTag struct {
	Name string `json:"name" validate:"required"`
}

func CreateHandler(svc Service, log *zap.SugaredLogger) handler {
	return handler{service: svc, log: log}
}

// CreateTags godoc
// @Summary      Create a tag
// @Description  Create a tag and return it
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param tag body InputTag true "A payload of new tag"
// @Success      201  {object}  web.RespBody{data=tagging.Response} "Response body for a new tag"
// @Failure      400  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      404  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /tags [post]
func (t handler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	payload := InputTag{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return bareknews.ErrInvalidJSON
	}

	tagRes, err := t.service.Create(ctx, payload.Name)
	if err != nil {
		if errors.Is(err, bareknews.ErrDataAlreadyExist) {
			return web.NewRequestError(bareknews.ErrDataAlreadyExist, http.StatusConflict)
		}
		return err
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfully creating a tag",
		Data:    tagRes,
	}

	return web.Respond(w, payloadRes, http.StatusCreated)
}

// GetTagById godoc
// @Summary      Get a tag
// @Description  Get a tag by id
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tag ID"  Format(uuid)
// @Success      200  {object}  web.RespBody{data=tagging.Response} "Response body for a tag"
// @Failure      404  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /tags/{id} [get]
func (t handler) GetById(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rawId := chi.URLParam(r, "tagId")

	id, err := uuid.Parse(rawId)
	if err != nil {
		return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
	}

	tg, err := t.service.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
		}
		return err
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfully getting a tag",
		Data:    tg,
	}

	return web.Respond(w, payloadRes, http.StatusOK)
}

// UpdateTags godoc
// @Summary      Update a tag
// @Description  Update a tag and return it
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tag ID"  Format(uuid)
// @Param tag body InputTag true "A payload of new tag"
// @Success      200  {object}  web.RespBody{data=tagging.Response} "Response body for a new tag"
// @Failure      400  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      404  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /tags/{id} [put]
func (t handler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rawId := chi.URLParam(r, "tagId")

	id, err := uuid.Parse(rawId)
	if err != nil {
		return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
	}

	payload := InputTag{}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return bareknews.ErrInvalidJSON
	}

	tg, err := t.service.Update(ctx, id, payload.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
		}
		return err
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfully updating a tag",
		Data:    tg,
	}

	return web.Respond(w, payloadRes, http.StatusOK)
}

// GetAllTags godoc
// @Summary      Get all tags
// @Description  Get all tags
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200  {object}  web.RespBody{data=[]tagging.Response} "Array of tag body"
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /tags [get]
func (t handler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	tgs, err := t.service.GetAll(ctx)
	if err != nil {
		return err
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfully getting all tags",
		Data:    tgs,
	}

	return web.Respond(w, payloadRes, http.StatusOK)
}

// DeleteTags godoc
// @Summary      Delete a tag
// @Description  Delete a tag by id
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tag ID"  Format(uuid)
// @Success      200  {object}  web.RespBody{data=object}
// @Failure      404  {object}  web.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  web.ErrRespBody{error=object{message=string}}
// @Router       /tags/{id} [delete]
func (t handler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rawId := chi.URLParam(r, "tagId")

	id, err := uuid.Parse(rawId)
	if err != nil {
		return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
	}

	err = t.service.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.NewRequestError(bareknews.ErrDataNotFound, http.StatusNotFound)
		}
		return err
	}

	payloadRes := web.GeneralResponse{
		Message: "Successfully deleting a tag",
		Data:    struct{}{},
	}

	return web.Respond(w, payloadRes, http.StatusOK)
}
