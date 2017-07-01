package alarmy_test

import (
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/alarmy"
)

var store alarmy.Store
var api *alarmy.Api
var boltStore *alarmy.BoltStore

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	boltStore.Close()
	os.Exit(retCode)
}

func setup() {
	var err error
	boltStore, err = alarmy.NewBoltStore("alarmy_test.db")
	if err != nil {
		panic(err)
	}

	store = alarmy.Store{
		ProjectStore: alarmy.NewBoltProjectStore(boltStore),
	}

	api = alarmy.NewApi(store)
}
