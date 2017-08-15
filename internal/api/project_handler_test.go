package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/jvikstedt/alarmy/internal/model"
)

func TestProjectCreate(t *testing.T) {
	testCases := []struct {
		name      string
		raw       string
		status    int
		illegalID int
		increment int
		errors    []model.Error
	}{
		{
			name:      "success",
			raw:       `{"id": 9999, "name": "Golang"}`,
			status:    http.StatusCreated,
			increment: 1,
			illegalID: 9999,
			errors:    []model.Error{},
		}, {
			name:   "validation",
			raw:    `{"name": ""}`,
			status: http.StatusUnprocessableEntity,
			errors: []model.Error{
				model.Error{Type: "validation", Name: "name", Reason: "missing"},
			},
		}, {
			name:   "invalid_input",
			raw:    "",
			status: http.StatusBadRequest,
			errors: []model.Error{
				model.Error{Type: "invalid_input", Name: "project", Reason: "EOF"},
			},
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			actualProject := model.Project{}
			json.Unmarshal([]byte(v.raw), &actualProject)

			b := bytes.NewBufferString(v.raw)
			var rr *httptest.ResponseRecorder

			IsIncrementedBy(t, "projects", v.increment, str.Project().CountAll, func() {
				rr, _ = MakeRequest("POST", "/projects", b)
			})

			if rr.Code != v.status {
				t.Fatalf("Expected status: %d but got %d", v.status, rr.Code)
			}

			response := model.ProjectResponse{}
			err := json.Unmarshal([]byte(rr.Body.String()), &response)
			if err != nil {
				t.Fatal(err)
			}

			if rr.Code == http.StatusCreated {
				str.Project().FindByID(int(response.Project.ID), &model.Project{})

				if response.Project.ID == 0 {
					t.Fatal("Created record id should not be 0")
				}

				if int(response.Project.ID) == v.illegalID {
					t.Fatal("User should not be able to set ID")
				}
			} else {
				if response.Project.ID != 0 {
					t.Fatalf("Created record id should be 0 was %d", response.Project.ID)
				}
			}

			if response.Project.Name != actualProject.Name {
				t.Fatalf("Expected Project.Name %s but got %s", actualProject.Name, response.Project.Name)
			}

			if same := reflect.DeepEqual(response.Errors, v.errors); !same {
				t.Fatalf("Expected errors %v got %v", v.errors, response.Errors)
			}
		})
	}
}

func TestProjectAll(t *testing.T) {
	str.Project().Clear()

	testCases := []struct {
		name    string
		project *model.Project
	}{
		{name: "nil", project: nil},
		{name: "golang", project: &model.Project{Name: "Golang"}},
		{name: "ruby", project: &model.Project{Name: "Ruby"}},
	}

	for i, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			if v.project != nil {
				str.Project().Create(v.project)
			}

			rr, _ := MakeRequest("GET", "/projects", nil)

			if rr.Code != http.StatusOK {
				t.Fatalf("Expected status: %d but got %d", http.StatusOK, rr.Code)
			}

			response := model.ProjectListResponse{}
			err := json.Unmarshal([]byte(rr.Body.String()), &response)
			if err != nil {
				t.Fatal(err)
			}

			if len(response.Projects) != i {
				t.Fatalf("Expected %d projects but got %d", i, len(response.Projects))
			}
		})
	}
}

// Test data
//var testProjects = []model.Project{
//	model.Project{ID: 1, Name: "Golang", CreatedAt: time.Now(), UpdatedAt: time.Now()},
//	model.Project{ID: 2, Name: "Ruby", CreatedAt: time.Now(), UpdatedAt: time.Now()},
//	model.Project{ID: 3, Name: "Javascript", CreatedAt: time.Now(), UpdatedAt: time.Now()},
//}
//
//func TestProjectAll(t *testing.T) {
//	handler, mockStore, _ := testDependencies()
//	req, _ := http.NewRequest("GET", "/projects", nil)
//
//	for i := 0; i < len(testProjects)+1; i++ {
//		rr := httptest.NewRecorder()
//
//		// Setup mockStore properly
//		mockStore.ProjectStore.Returns.Projects = testProjects[0:i]
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
//		// Go through projects and verify
//		for i, p := range projects {
//			assert.Equal(t, testProjects[i], p)
//		}
//	}
//
//	// Error case
//	mockStore.ProjectStore.Returns.Error = fmt.Errorf("Any error")
//	rr := httptest.NewRecorder()
//	handler.ServeHTTP(rr, req)
//	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
//}
//
//func TestProjectCreate(t *testing.T) {
//	handler, mockStore, _ := testDependencies()
//	b := bytes.NewBufferString(`{"Name": "Golang"}`)
//	req, _ := http.NewRequest("POST", "/projects", b)
//	req.Header.Add("Content-Type", "application/json")
//
//	rr := httptest.NewRecorder()
//
//	// Setup mock
//	testProject := model.Project{ID: 1, Name: "Golang", CreatedAt: time.Now(), UpdatedAt: time.Now()}
//	mockStore.ProjectStore.Returns.Project = testProject
//
//	handler.ServeHTTP(rr, req)
//
//	// Make sure mock receives correct object with correct information
//	assert.Equal(t, mockStore.ProjectStore.Receives.Project.Name, "Golang", "project name")
//
//	project := readProject(t, rr.Body)
//
//	assert.Equal(t, http.StatusCreated, rr.Code, "status code")
//	assert.Equal(t, testProject.Name, project.Name, "project name")
//	assert.Equal(t, testProject.ID, project.ID, "project id")
//
//	// Validation test
//	b = bytes.NewBufferString(`{"Name": ""}`)
//	req, _ = http.NewRequest("POST", "/projects", b)
//	req.Header.Add("Content-Type", "application/json")
//	rr = httptest.NewRecorder()
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code, "status code")
//	errors := readErrors(t, rr.Body)
//	assert.Contains(t, errors["name"], "is required", "name validation required")
//}
//
//func TestProjectGetOne(t *testing.T) {
//	handler, mockStore, _ := testDependencies()
//	for i, p := range testProjects {
//		req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d", i+1), nil)
//		rr := httptest.NewRecorder()
//
//		mockStore.ProjectStore.Returns.Project = p
//		handler.ServeHTTP(rr, req)
//
//		assert.Equal(t, http.StatusOK, rr.Code, "status code")
//		assert.Equal(t, p.ID, mockStore.ProjectStore.Receives.ID)
//		project := readProject(t, rr.Body)
//		assert.Equal(t, p.ID, project.ID, "project id")
//		assert.Equal(t, p.Name, project.Name, "project name")
//		assert.Equal(t, p.CreatedAt, project.CreatedAt, "project createdAt")
//		assert.Equal(t, p.UpdatedAt, project.UpdatedAt, "project updatedAt")
//	}
//}
//
//func TestProjectDestroy(t *testing.T) {
//	handler, mockStore, _ := testDependencies()
//	// Happy case
//	req, _ := http.NewRequest("DELETE", "/projects/1", nil)
//	rr := httptest.NewRecorder()
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code, "status code")
//	assert.Equal(t, 1, mockStore.ProjectStore.Receives.ID, "project id")
//
//	// invalid id provided
//	req, _ = http.NewRequest("DELETE", "/projects/asddsa", nil)
//	rr = httptest.NewRecorder()
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
//
//	// DB returns an error
//	mockStore.ProjectStore.Returns.Error = fmt.Errorf("any error")
//	rr = httptest.NewRecorder()
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status code")
//}
//
//// readProject is a helper function to read a single project from the body
//func readProject(t *testing.T, r io.Reader) model.Project {
//	data, err := ioutil.ReadAll(r)
//	if err != nil {
//		t.Errorf("Error reading from body %e", err)
//	}
//
//	project := model.Project{}
//	err = json.Unmarshal(data, &project)
//	if err != nil {
//		t.Errorf("JSON unmarshalling error %e", err)
//	}
//
//	return project
//}
//
//func TestProjectUpdate(t *testing.T) {
//	handler, mockStore, _ := testDependencies()
//	// Happy case
//	b := bytes.NewBufferString(`{"Name": "Test"}`)
//
//	req, _ := http.NewRequest("PATCH", "/projects/1", b)
//	req.Header.Add("Content-Type", "application/json")
//	rr := httptest.NewRecorder()
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusOK, rr.Code, "status code")
//	assert.Equal(t, "Test", mockStore.ProjectStore.Receives.Project.Name)
//	assert.Equal(t, 1, mockStore.ProjectStore.Receives.Project.ID)
//
//	// Invalid name
//	b = bytes.NewBufferString(`{"Name": ""}`)
//
//	req, _ = http.NewRequest("PATCH", "/projects/1", b)
//	req.Header.Add("Content-Type", "application/json")
//	rr = httptest.NewRecorder()
//	mockStore.ProjectStore.CallCount = 0
//	handler.ServeHTTP(rr, req)
//
//	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code, "status code")
//	assert.Equal(t, 0, mockStore.ProjectStore.CallCount, "store should not be called")
//}
//
//// readProjects is a helper function to read array of projects from the body
//func readProjects(t *testing.T, r io.Reader) []model.Project {
//	data, err := ioutil.ReadAll(r)
//	if err != nil {
//		t.Errorf("Error reading from body %e", err)
//	}
//
//	projects := []model.Project{}
//	err = json.Unmarshal(data, &projects)
//	if err != nil {
//		t.Errorf("JSON unmarshalling error %e", err)
//	}
//
//	return projects
//}
//
//// readErrors is a helper function to read map of errors from the body
//func readErrors(t *testing.T, r io.Reader) map[string][]string {
//	data, err := ioutil.ReadAll(r)
//	if err != nil {
//		t.Errorf("Error reading from body %e", err)
//	}
//
//	errors := make(map[string][]string)
//	err = json.Unmarshal(data, &errors)
//	if err != nil {
//		t.Errorf("JSON unmarshalling error %e", err)
//	}
//
//	return errors
//}
