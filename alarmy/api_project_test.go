package alarmy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProjectAll(t *testing.T) {
	req, _ := http.NewRequest("GET", "/projects", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.ProjectAll)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
