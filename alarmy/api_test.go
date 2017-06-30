package alarmy_test

import (
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/alarmy"
)

var store alarmy.Store
var api *alarmy.Api

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	store.Close()
	os.Exit(retCode)
}

func setup() {
	var err error
	store, err = alarmy.NewBoltStore("alarmy_test.db")
	if err != nil {
		panic(err)
	}

	api = alarmy.NewApi(store)
}
