package editor_test

import (
	"bytes"
	"io"
	"sync"
	"testing"

	edit "github.com/jvikstedt/alarmy/editor"
)

type Project struct {
	Name string
}

func TestNewObject(t *testing.T) {
	var resources = map[string]edit.Resource{
		"project": edit.Resource{
			Fields: []edit.Field{
				edit.Field{Name: "Name"},
			},
			New: func() interface{} { return &Project{} },
		},
	}

	pr, pw := io.Pipe()
	var b bytes.Buffer

	editor := edit.NewEditor(pr, &b, resources)
	_, err := editor.NewObject("something")

	if err == nil {
		t.Errorf("Expected error didn't occur when trying to get object with not specified resource key")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		object, err := editor.NewObject("project")
		if err != nil {
			t.Errorf("Error during NewObject (project) creation %v", err)
		}

		project, ok := object.(*Project)
		if !ok {
			t.Errorf("Expected Project object but did not get one")
		}
		Equal(t, "project name", "Golang", project.Name)
	}()

	pw.Write([]byte("Golang\n"))
	wg.Wait()
}

func Equal(t *testing.T, name string, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("%s\nExpected: %v\nWas: %v\n", name, expected, actual)
	}
}
