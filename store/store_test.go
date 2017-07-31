package store_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/store"
	_ "github.com/mattn/go-sqlite3"
)

var stores = make(map[string]store.Store)
var currentStore store.Store

func TestMain(m *testing.M) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStore := store.NewSqlStore(db, "sqlite3")
	err = sqlStore.SetupTables()
	if err != nil {
		panic(err)
	}

	stores["sql"] = sqlStore

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
