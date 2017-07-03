package api_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/api"
	"github.com/jvikstedt/alarmy/store"
)

var handler http.Handler
var mockStore *store.MockStore

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	os.Exit(retCode)
}

func setup() {
	mockStore = &store.MockStore{}
	store := store.Store{
		ProjectStore: mockStore,
	}
	api := api.NewApi(store)

	var err error
	handler, err = api.Handler()
	if err != nil {
		panic(err)
	}
}
