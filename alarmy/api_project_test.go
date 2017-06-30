package alarmy_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jvikstedt/alarmy/alarmy"
	"github.com/stretchr/testify/assert"
)

// TestProjectAll tests GET /projects
func TestProjectAll(t *testing.T) {
	// Reset tables
	store.RecreateAllTables()

	// Setup request
	req, _ := http.NewRequest("GET", "/projects", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.ProjectAll)

	// Test data
	testProjects := []alarmy.Project{
		alarmy.Project{Name: "Golang"},
		alarmy.Project{Name: "Ruby"},
		alarmy.Project{Name: "Javascript"},
	}

	for i := 0; i < len(testProjects)+1; i++ {
		// Skip first creation to test without projects
		if i != 0 {
			store.ProjectCreate(testProjects[i-1])
		}

		// Make a request
		handler.ServeHTTP(rr, req)

		// Read projects from body
		projects := readProjects(t, rr.Body)

		assert.Equal(t, http.StatusOK, rr.Code, "status code")
		assert.Equal(t, i, len(projects), "projects length")

		for i, p := range projects {
			assert.Equal(t, testProjects[i].Name, p.Name, "project name")
		}
	}
}

func readProjects(t *testing.T, r io.Reader) []alarmy.Project {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading from body %e", err)
	}

	projects := []alarmy.Project{}
	err = json.Unmarshal(data, &projects)
	if err != nil {
		t.Errorf("JSON unmarshalling error %e", err)
	}

	return projects
}
