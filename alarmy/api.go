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
		r.Get("/", a.GetProjects)
	})

	return r, nil
}

func (a *Api) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := a.store.Projects()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, projects)
}
