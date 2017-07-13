package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jvikstedt/alarmy/model"
)

type JobRequest struct {
	model.Job
	ProtectedID   interface{} `json:"id,omitempty"`
	OmitCreatedAt interface{} `json:"created_at,omitempty"`
	OmitUpdatedAt interface{} `json:"updated_at,omitempty"`
}

func (p *JobRequest) Bind(r *http.Request) error {
	return nil
}

// JobAll handler for getting all projects
func (a *Api) JobAll(w http.ResponseWriter, r *http.Request) {
	projects, err := a.store.JobAll()
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.JSON(w, r, projects)
}

// JobCreate handler for creating a project
func (a *Api) JobCreate(w http.ResponseWriter, r *http.Request) {
	data := &JobRequest{}
	err := render.Bind(r, data)
	if stop := a.CheckErr(w, r, err, http.StatusUnprocessableEntity); stop {
		return
	}

	if errors := data.Job.Errors(); len(errors) > 0 {
		a.Printf(r.Context(), "%v", errors)
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, errors)
		return
	}

	project, err := a.store.JobCreate(data.Job)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, project)
}

// JobGetOne handler to get single project by id
func (a *Api) JobGetOne(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	project, err := a.store.JobGetOne(projectID)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, project)
}

// JobDestroy delete a single project by id
func (a *Api) JobDestroy(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	err = a.store.JobDestroy(projectID)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
}

// JobUpdate update a project by id
func (a *Api) JobUpdate(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	data := &JobRequest{}
	err = render.Bind(r, data)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}
	data.Job.ID = projectID

	if errors := data.Job.Errors(); len(errors) > 0 {
		a.HandleError(w, r, errors, http.StatusUnprocessableEntity)
		render.JSON(w, r, errors)
		return
	}

	project, err := a.store.JobUpdate(data.Job)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, project)
}
