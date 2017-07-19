package store

import "github.com/jvikstedt/alarmy/model"

type Store interface {
	Project() ProjectStore
	Job() JobStore
}

type ProjectStore interface {
	All() ([]model.Project, error)
	Create(model.Project) (model.Project, error)
	Update(model.Project) (model.Project, error)
	Destroy(int) error
	GetOne(int) (model.Project, error)
	RemoveAll() error
}

type JobStore interface {
	All() ([]model.Job, error)
	Create(model.Job) (model.Job, error)
	Update(model.Job) (model.Job, error)
	Destroy(int) error
	GetOne(int) (model.Job, error)
	RemoveAll() error
}
