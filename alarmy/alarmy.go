package alarmy

import (
	"net/http"
	"time"
)

type Store interface {
	ProjectAll() ([]Project, error)
	ProjectCreate(Project) (Project, error)
}

type Router interface {
	Setup() (http.Handler, error)
}

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectRequest struct {
	Project
	ProtectedID   interface{} `json:"id,omitempty"`
	OmitCreatedAt interface{} `json:"created_at,omitempty"`
	OmitUpdatedAt interface{} `json:"updated_at,omitempty"`
}

func (a *ProjectRequest) Bind(r *http.Request) error {
	return nil
}
