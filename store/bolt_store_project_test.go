package store_test

import (
	"testing"
	"time"

	"github.com/jvikstedt/alarmy/model"
	"github.com/stretchr/testify/assert"
)

func TestProjectAll(t *testing.T) {
	currentStore.ProjectRemoveAll()

	// Zero projects
	projects, err := currentStore.ProjectAll()
	assert.Nil(t, err, "project all should not return an error")
	assert.Equal(t, 0, len(projects), "should be 0 projects")

	project := model.Project{Name: "Golang"}
	currentStore.ProjectCreate(project)

	// One project
	projects, err = currentStore.ProjectAll()
	assert.Nil(t, err, "project all should not return an error")
	assert.Equal(t, 1, len(projects), "should be 1 project")
}

func TestProjectCreate(t *testing.T) {
	currentStore.ProjectRemoveAll()

	testProjects := []model.Project{
		model.Project{Name: "Golang"},
		model.Project{Name: "Ruby"},
		model.Project{Name: "Javascript"},
	}

	for i, p := range testProjects {
		project, err := currentStore.ProjectCreate(p)

		assert.Nil(t, err, "ProjectCreate should not return an error")
		assert.Equal(t, i+1, project.ID, "id should be the same")
		assert.Equal(t, p.Name, project.Name, "project name should be same")
		duration := time.Since(project.CreatedAt).Seconds()
		assert.True(t, duration > 0 && duration < 1, "duration since created at should be between 0 and 1 seconds")
		duration = time.Since(project.UpdatedAt).Seconds()
		assert.True(t, duration > 0 && duration < 1, "duration since updated at should be between 0 and 1 seconds")
	}
}

func TestProjectUpdate(t *testing.T) {
	project, _ := currentStore.ProjectCreate(model.Project{Name: "Javascript"})
	project.Name = "Golang"

	testProject, err := currentStore.ProjectUpdate(project)
	assert.Nil(t, err, "ProjectUpdate should not return an error")
	assert.Equal(t, project.Name, testProject.Name, "project name should be the updated one")

	afterUpdateProject, _ := currentStore.ProjectGetOne(project.ID)
	assert.True(t, project.UpdatedAt.Before(afterUpdateProject.UpdatedAt), "UpdatedAt should be updated after project update")
	assert.Equal(t, project.Name, afterUpdateProject.Name, "project name should be the updated one")
}

func TestProjectDestroy(t *testing.T) {
	project, _ := currentStore.ProjectCreate(model.Project{Name: "Golang"})

	_, err := currentStore.ProjectGetOne(project.ID)
	assert.Nil(t, err, "project should be in the store before destroy")
	err = currentStore.ProjectDestroy(project.ID)
	assert.Nil(t, err, "ProjectDestroy should not return an error")
	_, err = currentStore.ProjectGetOne(project.ID)
	assert.Error(t, err, "project should not be in the store after destroy")
}

func TestProjectGetOne(t *testing.T) {
	initialProject, _ := currentStore.ProjectCreate(model.Project{Name: "Golang"})

	testProject, err := currentStore.ProjectGetOne(initialProject.ID)
	assert.Nil(t, err, "ProjectGetOne should not return an error")
	assert.Equal(t, initialProject.Name, testProject.Name, "ProjectGetOne should return correct project")
}
