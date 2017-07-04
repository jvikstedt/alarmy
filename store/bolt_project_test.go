package store_test

import (
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/model"
	"github.com/jvikstedt/alarmy/store"
	"github.com/stretchr/testify/assert"
)

var boltStore store.BoltStore
var boltProjectStore *store.BoltProjectStore

func TestMain(m *testing.M) {
	boltStore, err := store.NewBoltStore("alarmy_test.db")
	if err != nil {
		panic(err)
	}

	err = boltStore.RecreateAllTables()
	if err != nil {
		panic(err)
	}

	boltProjectStore = store.NewBoltProjectStore(boltStore)

	retCode := m.Run()
	boltStore.Close()
	os.Exit(retCode)
}

func BoltProjectAllTest(t *testing.T) {
	projects, err := boltProjectStore.ProjectAll()

	assert.Nil(t, err, "project all should not return an error")
	assert.Equal(t, 0, len(projects), "should be 0 projects")

	project := model.Project{Name: "Golang"}
	boltProjectStore.ProjectCreate(project)

	projects, err = boltProjectStore.ProjectAll()

	assert.Nil(t, err, "project all should not return an error")
	assert.Equal(t, 1, len(projects), "should be 1 project")
}
