package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/jvikstedt/alarmy/internal/model"
)

type ProjectRequest struct {
	*model.Project
	ProtectedID   interface{} `json:"id,omitempty"`
	OmitCreatedAt interface{} `json:"created_at,omitempty"`
	OmitUpdatedAt interface{} `json:"updated_at,omitempty"`
	OmitDeletedAt interface{} `json:"deleted_at,omitempty"`
}

func (p *ProjectRequest) Bind(r *http.Request) error {
	return nil
}

// ProjectAll handler for getting all projects
func (a *Api) ProjectAll(w http.ResponseWriter, r *http.Request) {
	projects := []model.Project{}
	if err := a.str.Project().GetAll(&projects); err != nil {
		a.Printf(r.Context(), "%v", err)
		render.Render(w, r, NewResponse(http.StatusInternalServerError, nil, nil))
	}

	render.Render(w, r, NewResponse(http.StatusOK, projects, nil))
}

// ProjectCreate handler for creating a project
func (a *Api) ProjectCreate(w http.ResponseWriter, r *http.Request) {
	data := &ProjectRequest{}
	if err := render.Bind(r, data); err != nil {
		errors := []model.Error{
			model.Error{Type: "invalid_input", Name: "project", Reason: err.Error()},
		}
		render.Render(w, r, NewResponse(http.StatusBadRequest, nil, errors))
		return
	}

	project := data.Project
	if errors := project.Errors(); len(errors) > 0 {
		render.Render(w, r, NewResponse(http.StatusUnprocessableEntity, nil, errors))
		return
	}
	err := a.str.Project().Create(project)
	if err != nil {
		render.Render(w, r, NewResponse(http.StatusInternalServerError, nil, nil))
		return
	}

	render.Render(w, r, NewResponse(http.StatusCreated, project, nil))
}

// ProjectGetOne handler to get single project by id
func (a *Api) ProjectGetOne(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if err != nil {
		errors := []model.Error{
			model.Error{Type: "parse_error", Name: "projectID", Reason: err.Error()},
		}
		render.Render(w, r, NewResponse(http.StatusBadRequest, nil, errors))
		return
	}

	project := model.Project{ID: uint(projectID)}
	err = a.str.Project().Find(&project)
	if err != nil {
		errors := []model.Error{
			model.Error{Type: "store_error", Name: "project", Reason: err.Error()},
		}
		render.Render(w, r, NewResponse(http.StatusNotFound, nil, errors))
		return
	}

	render.Render(w, r, NewResponse(http.StatusOK, project, nil))
}

// ProjectDestroy delete a single project by id
func (a *Api) ProjectDestroy(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if err != nil {
		errors := []model.Error{
			model.Error{Type: "parse_error", Name: "projectID", Reason: err.Error()},
		}
		render.Render(w, r, NewResponse(http.StatusBadRequest, nil, errors))
		return
	}

	project := model.Project{ID: uint(projectID)}
	err = a.str.Project().Delete(&project)
	if err != nil {
		errors := []model.Error{
			model.Error{Type: "store_error", Name: "project", Reason: err.Error()},
		}
		render.Render(w, r, NewResponse(http.StatusNotFound, nil, errors))
		return
	}

	render.Render(w, r, NewResponse(http.StatusOK, project, nil))
}

// ProjectUpdate update a project by id
func (a *Api) ProjectUpdate(w http.ResponseWriter, r *http.Request) {
	projectID, err := a.URLParamInt(r, "projectID")
	if err != nil {
		errors := []model.Error{
			model.Error{Type: "parse_error", Name: "projectID", Reason: err.Error()},
		}
		render.Render(w, r, NewResponse(http.StatusBadRequest, nil, errors))
		return
	}

	project := model.Project{ID: uint(projectID)}
	err = a.str.Project().Find(&project)
	if err != nil {
		errors := []model.Error{
			model.Error{Type: "store_error", Name: "project", Reason: err.Error()},
		}
		render.Render(w, r, NewResponse(http.StatusNotFound, nil, errors))
		return
	}

	data := &ProjectRequest{}
	data.Project = &project
	if err := render.Bind(r, data); err != nil {
		errors := []model.Error{
			model.Error{Type: "invalid_input", Name: "project", Reason: err.Error()},
		}
		render.Render(w, r, NewResponse(http.StatusBadRequest, nil, errors))
		return
	}

	if errors := project.Errors(); len(errors) > 0 {
		render.Render(w, r, NewResponse(http.StatusUnprocessableEntity, nil, errors))
		return
	}
	err = a.str.Project().Update(&project)
	if err != nil {
		render.Render(w, r, NewResponse(http.StatusInternalServerError, nil, nil))
		return
	}

	render.Render(w, r, NewResponse(http.StatusOK, project, nil))
}
