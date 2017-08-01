package store

import "github.com/jvikstedt/alarmy/model"

type MockStore struct {
	ProjectStore *ProjectMockStore
}

type ProjectMockStore struct {
	CallCount int
	Returns   struct {
		Projects []model.Project
		Project  model.Project
		Error    error
	}
	Receives struct {
		Project model.Project
		ID      int
	}
}

func NewMockStore() *MockStore {
	return &MockStore{
		ProjectStore: &ProjectMockStore{},
	}
}

func (s *MockStore) Project() ProjectStore {
	return s.ProjectStore
}

func (s *MockStore) Job() JobStore {
	return nil
}

func (s *ProjectMockStore) All() ([]model.Project, error) {
	return s.Returns.Projects, s.Returns.Error
}

func (s *ProjectMockStore) Create(project *model.Project) error {
	s.Receives.Project = *project
	*project = s.Returns.Project
	return s.Returns.Error
}

func (s *ProjectMockStore) Update(project *model.Project) error {
	s.CallCount++
	s.Receives.Project = *project
	*project = s.Returns.Project
	return s.Returns.Error
}

func (s *ProjectMockStore) Destroy(id int) error {
	s.Receives.ID = id
	return s.Returns.Error
}

func (s *ProjectMockStore) GetOne(id int) (model.Project, error) {
	s.Receives.ID = id
	return s.Returns.Project, s.Returns.Error
}

func (s *ProjectMockStore) RemoveAll() error {
	return nil
}
