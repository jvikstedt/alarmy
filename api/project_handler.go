package api

import (
	"net/http"

	"github.com/go-chi/render"
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
	projects, err := a.store.Project().All()
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.JSON(w, r, projects)
}

// ProjectCreate handler for creating a project
func (a *Api) ProjectCreate(w http.ResponseWriter, r *http.Request) {
	data := &ProjectRequest{}
	err := render.Bind(r, data)
	if stop := a.CheckErr(w, r, err, http.StatusUnprocessableEntity); stop {
		return
	}

	if errors := data.Project.Errors(); len(errors) > 0 {
		a.Printf(r.Context(), "%v", errors)
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, errors)
		return
	}

	project, err := a.store.Project().Create(data.Project)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, project)
}

// ProjectGetOne handler to get single project by id
func (a *Api) ProjectGetOne(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	project, err := a.store.Project().GetOne(projectID)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, project)
}

// ProjectDestroy delete a single project by id
func (a *Api) ProjectDestroy(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	err = a.store.Project().Destroy(projectID)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
}

// ProjectUpdate update a project by id
func (a *Api) ProjectUpdate(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	data := &ProjectRequest{}
	err = render.Bind(r, data)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}
	data.Project.ID = projectID

	if errors := data.Project.Errors(); len(errors) > 0 {
		a.HandleError(w, r, errors, http.StatusUnprocessableEntity)
		render.JSON(w, r, errors)
		return
	}

	project, err := a.store.Project().Update(data.Project)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, project)
}
