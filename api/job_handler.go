package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jvikstedt/alarmy/model"
	"github.com/jvikstedt/alarmy/schedule"
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

// JobAll handler for getting all jobs
func (a *Api) JobAll(w http.ResponseWriter, r *http.Request) {
	jobs, err := a.store.Job().All()
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.JSON(w, r, jobs)
}

// JobCreate handler for creating a job
func (a *Api) JobCreate(w http.ResponseWriter, r *http.Request) {
	data := &JobRequest{}
	err := render.Bind(r, data)
	if stop := a.CheckErr(w, r, err, http.StatusUnprocessableEntity); stop {
		return
	}

	// Validations
	errors := data.Job.Errors()

	_, err = a.store.Project().GetOne(data.Job.ProjectID)
	if err != nil {
		errors["project_id"] = append(errors["project_id"], err.Error())
	}

	// Validate spec
	if err := a.scheduler.ValidateSpec(data.Job.Spec); data.Job.Active && err != nil {
		errors["spec"] = append(errors["spec"], err.Error())
	}

	if len(errors) > 0 {
		a.Printf(r.Context(), "%v", errors)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		render.JSON(w, r, errors)
		return
	}

	job := data.Job
	err = a.store.Job().Create(&job)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	// Start scheduled job
	if job.Active {
		a.scheduler.AddEntry(schedule.EntryID(job.ID), job.Spec, a.executor.Execute)
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, job)
}

// JobGetOne handler to get single job by id
func (a *Api) JobGetOne(w http.ResponseWriter, r *http.Request) {
	jobID, err := a.URLParamInt(r, "jobID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	job, err := a.store.Job().GetOne(jobID)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, job)
}

// JobDestroy delete a single job by id
func (a *Api) JobDestroy(w http.ResponseWriter, r *http.Request) {
	jobID, err := a.URLParamInt(r, "jobID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	err = a.store.Job().Destroy(jobID)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
}

// JobUpdate update a job by id
func (a *Api) JobUpdate(w http.ResponseWriter, r *http.Request) {
	jobID, err := a.URLParamInt(r, "jobID")
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	data := &JobRequest{}
	err = render.Bind(r, data)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}
	data.Job.ID = jobID

	if errors := data.Job.Errors(); len(errors) > 0 {
		a.HandleError(w, r, errors, http.StatusUnprocessableEntity)
		render.JSON(w, r, errors)
		return
	}

	job := data.Job
	err = a.store.Job().Update(&job)
	if stop := a.CheckErr(w, r, err, http.StatusInternalServerError); stop {
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, job)
}
