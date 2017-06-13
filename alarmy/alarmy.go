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
