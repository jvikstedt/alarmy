package store

import "github.com/jvikstedt/alarmy/internal/model"

type Store interface {
	Project() ProjectStore
}

type ProjectStore interface {
	GetAll(projects *[]model.Project) error
	Create(project *model.Project) error
	Find(project *model.Project) error
	Delete(project *model.Project) error
	Update(project *model.Project) error
	CountAll() (int, error)
	Clear()
}
