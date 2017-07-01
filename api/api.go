package api

import (
	"net/http"

	"github.com/jvikstedt/alarmy/store"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
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
	})

	return r, nil
}
