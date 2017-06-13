package alarmy

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
)

type Api struct {
	store Store
}

func NewApi(store Store) *Api {
	return &Api{
		store: store,
	}
}

func (a *Api) Setup() (http.Handler, error) {
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

func (a *Api) ProjectAll(w http.ResponseWriter, r *http.Request) {
	projects, err := a.store.ProjectAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, projects)
}

func (a *Api) ProjectCreate(w http.ResponseWriter, r *http.Request) {
	data := &ProjectRequest{}
	if err := render.Bind(r, data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	project, err := a.store.ProjectCreate(data.Project)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, project)
}
