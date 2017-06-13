package alarmy

import (
	"net/http"
)

type Store interface {
	ProjectAll() ([]Project, error)
	ProjectCreate(Project) (Project, error)
	Close() error
	EnsureTablesExist() error
	RecreateAllTables() error
}

type Router interface {
	Setup() (http.Handler, error)
}
