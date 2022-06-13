package tags

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

type InputTag struct {
	Name string `json:"name" validate:"required"`
}

func (t Restapi) Route(r chi.Router) {
	r.Post("/", t.Create)
	r.Get("/{tagId}", t.GetById)
	r.Put("/{tagId}", t.Update)
	r.Get("/", t.GetAll)
	r.Delete("/{tagId}", t.Delete)
}

// CreateTags godoc
// @Summary      Create a tag
// @Description  Create a tag and return it
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param tag body InputTag true "A payload of new tag"
// @Success      201  {object}  restapi.RespBody{data=tagging.Response} "Response body for a new tag"
// @Failure      400  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      404  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /tags [post]
func (t Restapi) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := InputTag{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrInvalidJSON)
		if err != nil {
			log.Println("(error) tags.handler.create: ", err.Error())
		}
		return
	}

	tagRes, err := t.Service.Create(ctx, payload.Name)
	if err != nil {
		err = restapi.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.create: ", err.Error())
		}
		return
	}

	payloadRes := restapi.RespBody{
		Message: "Successfully creating a tag",
		Data: tagRes,
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.create: ", err.Error())
	}
}

// GetTagById godoc
// @Summary      Get a tag
// @Description  Get a tag by id
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tag ID"  Format(uuid)
// @Success      200  {object}  restapi.RespBody{data=tagging.Response} "Response body for a tag"
// @Failure      404  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /tags/{id} [get]
func (t Restapi) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) tags.handler.getById: ", err.Error())
		}
		return
	}

	tg, err := t.Service.GetById(ctx, id)
	if err != nil {
		err = restapi.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.getById: ", err.Error())
		}
		return
	}

	payloadRes := restapi.RespBody{
		Message: "Successfully getting a tag",
		Data: tg,
	}

	w.WriteHeader(http.StatusOK)
	
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.getById: ", err.Error())
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
// @Success      200  {object}  restapi.RespBody{data=tagging.Response} "Response body for a new tag"
// @Failure      400  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      404  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /tags/{id} [put]
func (t Restapi) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) tags.handler.update: ", err.Error())
		}
		return
	}

	payload := InputTag{}

	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrInvalidJSON)
		if err != nil {
			log.Println("(error) tags.handler.update: ", err.Error())
		}
		return
	}

	tg, err := t.Service.Update(ctx, id, payload.Name)
	if err != nil {
		err = restapi.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.update: ", err.Error())
		}
		return
	}

	payloadRes := restapi.RespBody{
		Message: "Successfully updating a tag",
		Data: tg,
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.update: ", err.Error())
	}
}

// GetAllTags godoc
// @Summary      Get all tags
// @Description  Get all tags
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200  {object}  restapi.RespBody{data=[]tagging.Response} "Array of tag body"
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /tags [get]
func (t Restapi) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tgs, err := t.Service.GetAll(ctx)
	if err != nil {
		err = restapi.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.getAll: ", err.Error())
		}
		return
	}

	payloadRes := restapi.RespBody{
		Message: "Successfully getting all tags",
		Data: tgs,
	}

	w.WriteHeader(http.StatusOK)
	
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.getAll: ", err.Error())
	}
}

// DeleteTags godoc
// @Summary      Delete a tag
// @Description  Delete a tag by id
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Tag ID"  Format(uuid)
// @Success      200  {object}  restapi.RespBody{data=object}
// @Failure      404  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Failure      500  {object}  restapi.ErrRespBody{error=object{message=string}}
// @Router       /tags/{id} [delete]
func (t Restapi) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rawId := chi.URLParam(r, "tagId")
	id, err := uuid.Parse(rawId)
	if err != nil {
		err = restapi.WriteErrResponse(w, bareknews.ErrDataNotFound)
		if err != nil {
			log.Println("(error) tags.handler.delete: ", err.Error())
		}
		return
	}

	err = t.Service.Delete(ctx, id)
	if err != nil {
		err = restapi.WriteErrResponse(w, err)
		if err != nil {
			log.Println("(error) tags.handler.delete: ", err.Error())
		}
		return
	}

	payloadRes := restapi.RespBody{
		Message: "Successfully deleting a tag",
		Data: struct{}{},
	}

	w.WriteHeader(http.StatusOK)
	
	err = json.NewEncoder(w).Encode(payloadRes)
	if err != nil {
		log.Println("(error) tags.handler.delete: ", err.Error())
	}
}