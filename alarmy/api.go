package alarmy

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

type Api struct {
}

func NewApi() *Api {
	return &Api{}
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
	w.Write([]byte("hi"))
}
