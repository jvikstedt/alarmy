package api

import (
	"net/http"

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

// ProjectAll handler for getting all projects
func (a *Api) ProjectAll(w http.ResponseWriter, r *http.Request) {
	projects, err := a.store.ProjectAll()
	if err != nil {
		a.Printf(r.Context(), "%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, projects)
}

// ProjectCreate handler for creating a project
func (a *Api) ProjectCreate(w http.ResponseWriter, r *http.Request) {
	data := &ProjectRequest{}
	if err := render.Bind(r, data); err != nil {
		a.Printf(r.Context(), "%v", err)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if errors := data.Project.Errors(); len(errors) > 0 {
		a.Printf(r.Context(), "%v", errors)
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, errors)
		return
	}

	project, err := a.store.ProjectCreate(data.Project)
	if err != nil {
		a.Printf(r.Context(), "%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, project)
}

// ProjectGetOne handler to get single project by id
func (a *Api) ProjectGetOne(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if err != nil {
		a.Printf(r.Context(), "%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	project, err := a.store.ProjectGetOne(projectID)

	if err != nil {
		a.Printf(r.Context(), "%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, project)
}

// ProjectDestroy delete a single project by id
func (a *Api) ProjectDestroy(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if err != nil {
		a.Printf(r.Context(), "%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = a.store.ProjectDestroy(projectID)
	if err != nil {
		a.Printf(r.Context(), "%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
}
