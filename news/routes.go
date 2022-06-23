package news

import (
	"github.com/Iiqbal2000/bareknews/pkg/web"
	"github.com/Iiqbal2000/bareknews/tags"
	"github.com/go-chi/chi/v5"
)

func Route(cfg web.APIMuxConfig, store Repository, tagSvc tags.Service) func(r chi.Router) {
	svc := CreateSvc(store, tagSvc)
	handler := CreateHandler(svc, cfg.Log)

	return func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/{newsId}", handler.GetById)
		r.Put("/{newsId}", handler.Update)
		r.Delete("/{newsId}", handler.Delete)
		r.Get("/", handler.GetAll)
	}
}
