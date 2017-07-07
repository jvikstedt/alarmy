package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jvikstedt/alarmy/store"
)

type Api struct {
	store store.Store
}

func NewApi(store store.Store) *Api {
	return &Api{
		store: store,
	}
}

func (a *Api) Handler() (http.Handler, error) {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Setup routes
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", a.ProjectAll)
		r.Post("/", a.ProjectCreate)
		r.Route("/{projectID}", func(r chi.Router) {
			r.Get("/", a.ProjectGetOne)
		})
	})

	return r, nil
}
