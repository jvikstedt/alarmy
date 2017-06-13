package alarmy

import (
	"net/http"

	"github.com/pressly/chi/render"
)

func (a *Api) ProjectAll(w http.ResponseWriter, r *http.Request) {
	projects, err := a.store.ProjectAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if projects == nil {
		projects = []Project{}
	}

	render.JSON(w, r, projects)
}

func (a *Api) ProjectCreate(w http.ResponseWriter, r *http.Request) {
	data := &ProjectRequest{}
	if err := render.Bind(r, data); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
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
