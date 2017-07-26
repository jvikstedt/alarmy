package store_test

import (
	"testing"
	"time"

	"github.com/jvikstedt/alarmy/model"
	"github.com/stretchr/testify/assert"
)

func TestProjectAll(t *testing.T) {
	currentStore.Project().RemoveAll()

	// Zero projects
	projects, err := currentStore.Project().All()
	assert.Nil(t, err, "project all should not return an error")
	assert.Equal(t, 0, len(projects), "should be 0 projects")

	project := model.Project{Name: "Golang"}
	currentStore.Project().Create(project)

	// One project
	projects, err = currentStore.Project().All()
	assert.Nil(t, err, "project all should not return an error")
	assert.Equal(t, 1, len(projects), "should be 1 project")
}

func TestProjectCreate(t *testing.T) {
	currentStore.Project().RemoveAll()

	testProjects := []model.Project{
		model.Project{Name: "Golang"},
		model.Project{Name: "Ruby"},
		model.Project{Name: "Javascript"},
	}

	for i, p := range testProjects {
		project, err := currentStore.Project().Create(p)

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
	project, _ := currentStore.Project().Create(model.Project{Name: "Javascript"})
	project.Name = "Golang"

	testProject, err := currentStore.Project().Update(project)
	assert.Nil(t, err, "ProjectUpdate should not return an error")
	assert.Equal(t, project.Name, testProject.Name, "project name should be the updated one")

	afterUpdateProject, _ := currentStore.Project().GetOne(project.ID)
	assert.True(t, project.UpdatedAt.Before(afterUpdateProject.UpdatedAt), "UpdatedAt should be updated after project update")
	assert.Equal(t, project.Name, afterUpdateProject.Name, "project name should be the updated one")
}

func TestProjectDestroy(t *testing.T) {
	project, _ := currentStore.Project().Create(model.Project{Name: "Golang"})

	_, err := currentStore.Project().GetOne(project.ID)
	assert.Nil(t, err, "project should be in the store before destroy")
	err = currentStore.Project().Destroy(project.ID)
	assert.Nil(t, err, "ProjectDestroy should not return an error")
	_, err = currentStore.Project().GetOne(project.ID)
	assert.Error(t, err, "project should not be in the store after destroy")
}

func TestProjectGetOne(t *testing.T) {
	initialProject, _ := currentStore.Project().Create(model.Project{Name: "Golang"})

	testProject, err := currentStore.Project().GetOne(initialProject.ID)
	assert.Nil(t, err, "ProjectGetOne should not return an error")
	assert.Equal(t, initialProject.Name, testProject.Name, "ProjectGetOne should return correct project")
}
