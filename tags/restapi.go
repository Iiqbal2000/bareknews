package tags

import (
	"encoding/json"
	"net/http"

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
func (t handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := InputTag{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, bareknews.ErrInvalidJSON)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	tagRes, err := t.service.Create(ctx, payload.Name)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, err)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadRes := web.RespBody{
		Message: "Successfully creating a tag",
		Data:    tagRes,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		t.log.Error(errors.Wrap(err, "failed to write a response"))
	}
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
func (t handler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, bareknews.ErrDataNotFound)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	tg, err := t.service.GetById(ctx, id)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, err)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadRes := web.RespBody{
		Message: "Successfully getting a tag",
		Data:    tg,
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		t.log.Error(errors.Wrap(err, "failed to write a response"))
	}
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
func (t handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, bareknews.ErrDataNotFound)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payload := InputTag{}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, bareknews.ErrInvalidJSON)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	tg, err := t.service.Update(ctx, id, payload.Name)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, err)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadRes := web.RespBody{
		Message: "Successfully updating a tag",
		Data:    tg,
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		t.log.Error(errors.Wrap(err, "failed to write a response"))
	}
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
func (t handler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tgs, err := t.service.GetAll(ctx)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, err)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadRes := web.RespBody{
		Message: "Successfully getting all tags",
		Data:    tgs,
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		t.log.Error(errors.Wrap(err, "failed to write a response"))
	}
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
func (t handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, bareknews.ErrDataNotFound)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	err = t.service.Delete(ctx, id)
	if err != nil {
		err = web.WriteErrResponse(w, t.log, err)
		if err != nil {
			t.log.Error(errors.Wrap(err, "failed to write a response"))
		}
		return
	}

	payloadRes := web.RespBody{
		Message: "Successfully deleting a tag",
		Data:    struct{}{},
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		t.log.Error(errors.Wrap(err, "failed to write a response"))
	}
}
