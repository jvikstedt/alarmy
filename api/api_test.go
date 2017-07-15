package api_test

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/api"
	"github.com/jvikstedt/alarmy/schedule"
	"github.com/jvikstedt/alarmy/store"
)

var handler http.Handler
var mockStore *store.MockStore
var logs = &bytes.Buffer{}

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

	mockScheduler := &schedule.MockScheduler{}

	logger := log.New(logs, "", log.LstdFlags)
	api := api.NewApi(store, logger, mockScheduler)

	var err error
	handler, err = api.Handler()
	if err != nil {
		panic(err)
	}
}
