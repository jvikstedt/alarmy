package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/jvikstedt/alarmy/internal"
	"github.com/jvikstedt/alarmy/internal/model"
	"github.com/jvikstedt/alarmy/internal/store"
	"github.com/jvikstedt/alarmy/schedule"
)

type Api struct {
	str       store.Store
	logger    *log.Logger
	scheduler schedule.Scheduler
	executor  *internal.Executor
}

func NewApi(str store.Store, logger *log.Logger, scheduler schedule.Scheduler, executor *internal.Executor) *Api {
	return &Api{
		str:       str,
		logger:    logger,
		scheduler: scheduler,
		executor:  executor,
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

	r.Route("/triggers", func(r chi.Router) {
		r.Get("/", a.TriggerAll)
		r.Post("/", a.TriggerCreate)
		r.Route("/{triggerID}", func(r chi.Router) {
			r.Get("/", a.TriggerGetOne)
			r.Delete("/", a.TriggerDestroy)
			r.Patch("/", a.TriggerUpdate)
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

type Response model.Response

func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewResponse(status int, data interface{}, errors []model.Error) *Response {
	if errors == nil {
		errors = []model.Error{}
	}
	return &Response{
		HTTPStatusCode: status,
		Data:           data,
		HasErrors:      len(errors) > 0,
		Errors:         errors,
	}
}
