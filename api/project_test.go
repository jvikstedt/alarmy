package api_test

import (
	"bytes"
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

func TestProjectCreate(t *testing.T) {
	b := bytes.NewBufferString(`{"Name": "Golang"}`)
	req, _ := http.NewRequest("POST", "/projects", b)
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Setup mock
	testProject := model.Project{ID: 1, Name: "Golang", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	mockStore.Project.ProjectCreate.Returns.Project = testProject

	handler.ServeHTTP(rr, req)

	// Make sure mock receives correct object with correct information
	assert.Equal(t, mockStore.Project.ProjectCreate.Receives.Project.Name, "Golang", "project name")

	project := readProject(t, rr.Body)

	assert.Equal(t, http.StatusCreated, rr.Code, "status code")
	assert.Equal(t, testProject.Name, project.Name, "project name")
	assert.Equal(t, testProject.ID, project.ID, "project id")

	// Validation test
	b = bytes.NewBufferString(`{"Name": ""}`)
	req, _ = http.NewRequest("POST", "/projects", b)
	req.Header.Add("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code, "status code")
	errors := readErrors(t, rr.Body)
	assert.Contains(t, errors["name"], "is required", "name validation required")
}

func TestProjectGetOne(t *testing.T) {
	for i, p := range testProjects {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d", i+1), nil)
		rr := httptest.NewRecorder()

		mockStore.Project.ProjectGetOne.Returns.Project = p
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code, "status code")
		assert.Equal(t, p.ID, mockStore.Project.ProjectGetOne.Receives.ID)
		project := readProject(t, rr.Body)
		assert.Equal(t, p.ID, project.ID, "project id")
		assert.Equal(t, p.Name, project.Name, "project name")
		assert.Equal(t, p.CreatedAt, project.CreatedAt, "project createdAt")
		assert.Equal(t, p.UpdatedAt, project.UpdatedAt, "project updatedAt")
	}
}

func TestProjectDestroy(t *testing.T) {
	// Happy case
	req, _ := http.NewRequest("DELETE", "/projects/1", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "status code")
	assert.Equal(t, 1, mockStore.Project.ProjectDestroy.Receives.ID, "project id")

	// invalid id provided
	req, _ = http.NewRequest("DELETE", "/projects/asddsa", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")

	// DB returns an error
	mockStore.Project.ProjectDestroy.Returns.Error = fmt.Errorf("any error")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
}

// readProject is a helper function to read a single project from the body
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

// readProjects is a helper function to read array of projects from the body
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

// readErrors is a helper function to read map of errors from the body
func readErrors(t *testing.T, r io.Reader) map[string][]string {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("Error reading from body %e", err)
	}

	errors := make(map[string][]string)
	err = json.Unmarshal(data, &errors)
	if err != nil {
		t.Errorf("JSON unmarshalling error %e", err)
	}

	return errors
}
