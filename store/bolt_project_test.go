package store_test

import (
	"testing"

	"github.com/jvikstedt/alarmy/model"
	"github.com/stretchr/testify/assert"
)

func BoltProjectAllTest(t *testing.T) {
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
