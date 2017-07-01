package alarmy

import (
	"net/http"
)

type Store struct {
	ProjectStore
}

type ProjectStore interface {
	ProjectAll() ([]Project, error)
	ProjectCreate(Project) (Project, error)
	ProjectUpdate(Project) (Project, error)
	ProjectDestroy(int) error
	ProjectGetOne(int) (Project, error)
	ProjectRemoveAll() error
}

type Router interface {
	Setup() (http.Handler, error)
}
