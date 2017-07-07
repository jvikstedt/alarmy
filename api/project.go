package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/render"
	"github.com/jvikstedt/alarmy/model"
)

type ProjectRequest struct {
	model.Project
	ProtectedID   interface{} `json:"id,omitempty"`
	OmitCreatedAt interface{} `json:"created_at,omitempty"`
	OmitUpdatedAt interface{} `json:"updated_at,omitempty"`
}

func (p *ProjectRequest) Bind(r *http.Request) error {
	return nil
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

func (a *Api) ProjectGetOne(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	if idStr := chi.URLParam(r, "projectID"); idStr != "" {
		projectID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		project, err = a.store.ProjectGetOne(projectID)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, project)
}
