package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

			IsChangedByAmount(t, "projects", v.increment, str.Project().CountAll, func() {
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
				str.Project().Find(&model.Project{ID: uint(response.Project.ID)})

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

func TestProjectGetOne(t *testing.T) {
	str.Project().Clear()

	testProject := model.Project{Name: "Golang"}
	err := str.Project().Create(&testProject)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("/projects/%d", testProject.ID), nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status: %d but got %d", http.StatusOK, rr.Code)
	}

	response := model.ProjectResponse{}
	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Project.Name != testProject.Name {
		t.Fatalf("Expected Project.Name to be %s but was %s", testProject.Name, response.Project.Name)
	}
}

func TestProjectDestroy(t *testing.T) {
	testProject := model.Project{Name: "Golang"}
	err := str.Project().Create(&testProject)
	if err != nil {
		t.Fatal(err)
	}
	testProject2 := model.Project{Name: "Ruby"}
	err = str.Project().Create(&testProject2)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/projects/%d", testProject.ID), nil)
	rr := httptest.NewRecorder()

	// Make sure only one record is deleted
	IsChangedByAmount(t, "projects", -1, str.Project().CountAll, func() {
		handler.ServeHTTP(rr, req)
	})

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status: %d but got %d", http.StatusOK, rr.Code)
	}

	response := model.ProjectResponse{}
	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Project.Name != testProject.Name {
		t.Fatalf("Expected Project.Name to be %s but was %s", testProject.Name, response.Project.Name)
	}
}

func TestProjectUpdate(t *testing.T) {
	testProject := model.Project{Name: "Golang"}
	err := str.Project().Create(&testProject)
	if err != nil {
		t.Fatal(err)
	}

	// Happy case
	b := bytes.NewBufferString(`{"Name": "Test"}`)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/projects/%d", testProject.ID), b)
	req.Header.Add("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status: %d but got %d", http.StatusOK, rr.Code)
	}

	response := model.ProjectResponse{}
	err = json.Unmarshal([]byte(rr.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Project.Name != "Test" {
		t.Fatalf("Expected Project.Name to be %s but was %s", "Test", response.Project.Name)
	}

	str.Project().Find(&testProject)

	if testProject.Name != "Test" {
		t.Fatalf("Expected Project.Name to be %s but was %s", "Test", testProject.Name)
	}
}
