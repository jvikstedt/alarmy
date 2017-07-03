package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jvikstedt/alarmy/model"
	"github.com/stretchr/testify/assert"
)

// Test data
var testProjects = []model.Project{
	model.Project{ID: 1, Name: "Golang", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	model.Project{ID: 2, Name: "Ruby", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	model.Project{ID: 3, Name: "Javascript", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

func TestProjectAll(t *testing.T) {
	req, _ := http.NewRequest("GET", "/projects", nil)

	for i := 0; i < len(testProjects)+1; i++ {
		rr := httptest.NewRecorder()

		// Setup mockStore properly
		mockStore.Project.ProjectAll.Returns.Projects = testProjects[0:i]

		// Make a request
		handler.ServeHTTP(rr, req)

		// Read projects from body
		projects := readProjects(t, rr.Body)

		assert.Equal(t, http.StatusOK, rr.Code, "status code")
		assert.Equal(t, i, len(projects), "projects length")

		// Go through projects and verify fields
		for i, p := range projects {
			assert.Equal(t, testProjects[i].ID, p.ID, "project id")
			assert.Equal(t, testProjects[i].Name, p.Name, "project name")
			assert.Equal(t, testProjects[i].CreatedAt, p.CreatedAt, "project CreatedAt")
			assert.Equal(t, testProjects[i].UpdatedAt, p.UpdatedAt, "project UpdatedAt")
		}
	}

	// Error case
	mockStore.Project.ProjectAll.Returns.Error = fmt.Errorf("Any error")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
}

//
//// TestProjectAll tests GET /projects
//func TestProjectAll(t *testing.T) {
//	// Reset table
//	store.ProjectRemoveAll()
//
//	// Setup request
//	req, _ := http.NewRequest("GET", "/projects", nil)
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(api.ProjectAll)
//
//	for i := 0; i < len(testProjects)+1; i++ {
//		// Skip first creation to test without projects
//		if i != 0 {
//			store.ProjectCreate(testProjects[i-1])
//		}
//
//		// Make a request
//		handler.ServeHTTP(rr, req)
//
//		// Read projects from body
//		projects := readProjects(t, rr.Body)
//
//		assert.Equal(t, http.StatusOK, rr.Code, "status code")
//		assert.Equal(t, i, len(projects), "projects length")
//
//		for i, p := range projects {
//			assert.Equal(t, testProjects[i].Name, p.Name, "project name")
//		}
//	}
//}
//
//func TestProjectCreate(t *testing.T) {
//	// Reset table
//	store.ProjectRemoveAll()
//
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(api.ProjectCreate)
//
//	for i, p := range testProjects {
//		b := new(bytes.Buffer)
//		json.NewEncoder(b).Encode(p)
//
//		req, _ := http.NewRequest("POST", "/projects", b)
//
//		// Make a request
//		handler.ServeHTTP(rr, req)
//
//		// Read project from body
//		project := readProject(t, rr.Body)
//
//		assert.Equal(t, http.StatusCreated, rr.Code, "status code")
//		assert.Equal(t, p.Name, project.Name, "project name")
//		assert.Equal(t, i+1, project.ID, "project id")
//	}
//
//	// Make sure data actually went to the database
//	projects, _ := store.ProjectAll()
//	assert.Equal(t, len(testProjects), len(projects), "projects length")
//	for i, p := range projects {
//		assert.Equal(t, testProjects[i].Name, p.Name, "project name")
//	}
//}
//

func readProject(t *testing.T, r io.Reader) model.Project {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading from body %e", err)
	}

	project := model.Project{}
	err = json.Unmarshal(data, &project)
	if err != nil {
		t.Errorf("JSON unmarshalling error %e", err)
	}

	return project
}

func readProjects(t *testing.T, r io.Reader) []model.Project {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading from body %e", err)
	}

	projects := []model.Project{}
	err = json.Unmarshal(data, &projects)
	if err != nil {
		t.Errorf("JSON unmarshalling error %e", err)
	}

	return projects
}
