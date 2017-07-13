package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: a.logger}))
	r.Use(middleware.Recoverer)

	// Setup routes
	r.Route("/projects", func(r chi.Router) {
		r.Get("/", a.ProjectAll)
		r.Post("/", a.ProjectCreate)
		r.Route("/{projectID}", func(r chi.Router) {
			r.Get("/", a.ProjectGetOne)
			r.Delete("/", a.ProjectDestroy)
			r.Patch("/", a.ProjectUpdate)
		})
	})

	r.Route("/jobs", func(r chi.Router) {
		r.Get("/", a.JobAll)
		r.Post("/", a.JobCreate)
		r.Route("/{jobID}", func(r chi.Router) {
			r.Get("/", a.JobGetOne)
			r.Delete("/", a.JobDestroy)
			r.Patch("/", a.JobUpdate)
		})
	})

	return r, nil
}

func (a *Api) Printf(ctx context.Context, format string, v ...interface{}) {
	a.logger.Printf("[%s] "+format, middleware.GetReqID(ctx), v)
}

func (a *Api) URLParamInt(r *http.Request, key string) (int, error) {
	asStr := chi.URLParam(r, key)
	if asStr == "" {
		return 0, fmt.Errorf("%s not set", key)
	}

	return strconv.Atoi(asStr)
}

func (a *Api) HandleError(w http.ResponseWriter, r *http.Request, err interface{}, statusCode int) {
	a.Printf(r.Context(), "%v", err)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (a *Api) CheckErr(w http.ResponseWriter, r *http.Request, err error, statusCode int) bool {
	if err != nil {
		a.HandleError(w, r, err, statusCode)
		return true
	}
	return false
}
