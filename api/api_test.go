package api_test

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jvikstedt/alarmy/alarm"
	"github.com/jvikstedt/alarmy/api"
	"github.com/jvikstedt/alarmy/schedule"
	"github.com/jvikstedt/alarmy/store"
)

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}

func testDependencies() (http.Handler, *store.MockStore, *bytes.Buffer) {
	mockStore := store.NewMockStore()
	mockScheduler := &schedule.MockScheduler{}

	logs := &bytes.Buffer{}
	logger := log.New(logs, "", log.LstdFlags)

	executor := alarm.NewExecutor(mockStore, logger)
	api := api.NewApi(mockStore, logger, mockScheduler, executor)

	var err error
	handler, err := api.Handler()
	if err != nil {
		panic(err)
	}

	return handler, mockStore, logs
}
