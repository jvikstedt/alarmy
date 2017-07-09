package store

import "github.com/jvikstedt/alarmy/model"

type MockStore struct {
	Project struct {
		ProjectAll struct {
			Returns struct {
				Projects []model.Project
				Error    error
			}
		}

		ProjectCreate struct {
			Receives struct {
				Project model.Project
			}
			Returns struct {
				Project model.Project
				Error   error
			}
		}

		ProjectGetOne struct {
			Receives struct {
				ID int
			}
			Returns struct {
				Project model.Project
				Error   error
			}
		}

		ProjectDestroy struct {
			Receives struct {
				ID int
			}
			Returns struct {
				Error error
			}
		}
	}
}

func (s *MockStore) ProjectAll() ([]model.Project, error) {
	return s.Project.ProjectAll.Returns.Projects, s.Project.ProjectAll.Returns.Error
}

func (s *MockStore) ProjectCreate(project model.Project) (model.Project, error) {
	s.Project.ProjectCreate.Receives.Project = project
	return s.Project.ProjectCreate.Returns.Project, s.Project.ProjectCreate.Returns.Error
}

func (s *MockStore) ProjectUpdate(project model.Project) (model.Project, error) {
	return model.Project{}, nil
}

func (s *MockStore) ProjectDestroy(id int) error {
	s.Project.ProjectDestroy.Receives.ID = id
	return s.Project.ProjectDestroy.Returns.Error
}

func (s *MockStore) ProjectGetOne(id int) (model.Project, error) {
	s.Project.ProjectGetOne.Receives.ID = id
	return s.Project.ProjectGetOne.Returns.Project, s.Project.ProjectGetOne.Returns.Error
}

func (s *MockStore) ProjectRemoveAll() error {
	return nil
}
