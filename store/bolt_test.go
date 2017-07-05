package store_test

import (
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/store"
)

var boltStore store.BoltStore
var projectStore store.ProjectStore

func TestMain(m *testing.M) {
	boltStore, err := store.NewBoltStore("alarmy_test.db")
	if err != nil {
		panic(err)
	}

	err = boltStore.RecreateAllTables()
	if err != nil {
		panic(err)
	}

	projectStore = store.NewBoltProjectStore(boltStore)

	retCode := m.Run()
	boltStore.Close()
	os.Exit(retCode)
}
