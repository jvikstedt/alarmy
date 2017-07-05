package store_test

import (
	"testing"
	"time"

	"github.com/jvikstedt/alarmy/model"
	"github.com/stretchr/testify/assert"
)

func TestProjectAll(t *testing.T) {
	projectStore.ProjectRemoveAll()

	// Zero projects
	projects, err := projectStore.ProjectAll()
	assert.Nil(t, err, "project all should not return an error")
	assert.Equal(t, 0, len(projects), "should be 0 projects")

	project := model.Project{Name: "Golang"}
	projectStore.ProjectCreate(project)

	// One project
	projects, err = projectStore.ProjectAll()
	assert.Nil(t, err, "project all should not return an error")
	assert.Equal(t, 1, len(projects), "should be 1 project")
}

func TestProjectCreate(t *testing.T) {
	projectStore.ProjectRemoveAll()

	testProjects := []model.Project{
		model.Project{Name: "Golang"},
		model.Project{Name: "Ruby"},
		model.Project{Name: "Javascript"},
	}

	for i, p := range testProjects {
		project, err := projectStore.ProjectCreate(p)

		assert.Nil(t, err, "ProjectCreate should not return an error")
		assert.Equal(t, i+1, project.ID, "id should be the same")
		assert.Equal(t, p.Name, project.Name, "project name should be same")
		duration := time.Since(project.CreatedAt).Seconds()
		assert.True(t, duration > 0 && duration < 1, "duration since created at should be between 0 and 1 seconds")
		duration = time.Since(project.UpdatedAt).Seconds()
		assert.True(t, duration > 0 && duration < 1, "duration since updated at should be between 0 and 1 seconds")
	}
}
