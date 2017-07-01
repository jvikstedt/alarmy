package store

import "github.com/jvikstedt/alarmy/model"

type Store struct {
	ProjectStore
}

type ProjectStore interface {
	ProjectAll() ([]model.Project, error)
	ProjectCreate(model.Project) (model.Project, error)
	ProjectUpdate(model.Project) (model.Project, error)
	ProjectDestroy(int) error
	ProjectGetOne(int) (model.Project, error)
	ProjectRemoveAll() error
}
