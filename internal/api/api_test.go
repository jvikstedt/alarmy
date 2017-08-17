package api_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jvikstedt/alarmy/internal"
	"github.com/jvikstedt/alarmy/internal/api"
	"github.com/jvikstedt/alarmy/internal/model"
	"github.com/jvikstedt/alarmy/internal/store"
	"github.com/jvikstedt/alarmy/schedule"
)

var handler http.Handler
var str store.Store

func TestMain(m *testing.M) {
	db, err := gorm.Open("sqlite3", "alarmy_test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(
		&model.Project{},
		&model.Job{},
		&model.Trigger{},
	)

	mockScheduler := &schedule.MockScheduler{}
	logs := &bytes.Buffer{}
	logger := log.New(logs, "", log.LstdFlags)

	str = store.NewGormStore(db)

	executor := internal.NewExecutor(db, logger)
	api := api.NewApi(str, logger, mockScheduler, executor)

	handler, err = api.Handler()
	if err != nil {
		panic(err)
	}

	retCode := m.Run()
	os.Exit(retCode)
}

func IsChangedByAmount(t *testing.T, message string, amount int, count func() (int, error), callback func()) {
	before, err := count()
	if err != nil {
		t.Fatal(err)
	}
	callback()
	after, err := count()
	if err != nil {
		t.Fatal(err)
	}
	if (after - before) != amount {
		t.Logf("%s, expected to increase by %d but was %d", message, amount, after-before)
	}
}

func MakeRequest(method string, url string, body io.Reader) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr, nil
}
