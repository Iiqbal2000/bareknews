package tags

import (
	"github.com/Iiqbal2000/bareknews/pkg/web"
	"github.com/go-chi/chi/v5"
)

type Routes struct {
	Router func(r chi.Router)
	Svc      Service
}

func Route(cfg web.APIMuxConfig, store Repository) Routes {
	svc := CreateSvc(store)
	handler := CreateHandler(svc, cfg.Log)

	r := Routes{
		Router: func(r chi.Router) {
			r.Post("/", handler.Create)
			r.Get("/{tagId}", handler.GetById)
			r.Put("/{tagId}", handler.Update)
			r.Get("/", handler.GetAll)
			r.Delete("/{tagId}", handler.Delete)
		},
		Svc: svc,
	}
	
	return r
}
