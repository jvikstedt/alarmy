package api

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jvikstedt/alarmy/store"
)

type Api struct {
	store  store.Store
	logger *log.Logger
}

func NewApi(store store.Store, logger *log.Logger) *Api {
	return &Api{
		store:  store,
		logger: logger,
	}
}

func (a *Api) Handler() (http.Handler, error) {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags)}))
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

func (a *Api) Printf(ctx context.Context, format string, v ...interface{}) {
	a.logger.Printf("[%s] "+format, middleware.GetReqID(ctx), v)
}
