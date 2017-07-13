package store

import "github.com/jvikstedt/alarmy/model"

type Store struct {
	ProjectStore
	JobStore
}

type ProjectStore interface {
	ProjectAll() ([]model.Project, error)
	ProjectCreate(model.Project) (model.Project, error)
	ProjectUpdate(model.Project) (model.Project, error)
	ProjectDestroy(int) error
	ProjectGetOne(int) (model.Project, error)
	ProjectRemoveAll() error
}

type JobStore interface {
	JobAll() ([]model.Job, error)
	JobCreate(model.Job) (model.Job, error)
	JobUpdate(model.Job) (model.Job, error)
	JobDestroy(int) error
	JobGetOne(int) (model.Job, error)
	JobRemoveAll() error
}
