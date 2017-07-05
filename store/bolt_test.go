package store_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/store"
)

var stores = make(map[string]store.Store)
var currentStore store.Store

func TestMain(m *testing.M) {
	boltStore, err := store.NewBoltStore("alarmy_test.db")
	if err != nil {
		panic(err)
	}
	defer boltStore.Close()

	err = boltStore.RecreateAllTables()
	if err != nil {
		panic(err)
	}

	stores["bolt"] = boltStore.Store()

	var result int
	for k, v := range stores {
		fmt.Printf("Setting store: %s\n", k)
		currentStore = v

		retCode := m.Run()
		if retCode != 0 {
			result = retCode
		}
	}

	os.Exit(result)
}