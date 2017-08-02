package edit

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Field is one of the Object fields
type Field struct {
	Name        string
	Association string
	Default     interface{}
}

// Resource is a instruction for the editor
// New is a function that tells how to make a new Resource
type Resource struct {
	Fields []Field
	New    func() interface{}
}

type Editor struct {
	reader    *bufio.Reader
	writer    io.Writer
	resources map[string]Resource
}

func NewEditor(reader io.Reader, writer io.Writer, resources map[string]Resource) *Editor {
	return &Editor{
		reader:    bufio.NewReader(reader),
		writer:    writer,
		resources: resources,
	}
}

func (e *Editor) NewObject(name string) (interface{}, error) {
	// Get Resource
	r, ok := e.resources[name]
	if !ok {
		return nil, fmt.Errorf("Could not find resource by name %s", name)
	}

	// Create a new object
	object := r.New()

	// Loop fields
	for _, f := range r.Fields {
		if err := e.handleField(object, f); err != nil {
			return nil, err
		}
	}

	return object, nil
}

func (e *Editor) UpdateObject(object interface{}, name string) error {
	// Get Resource
	r, ok := e.resources[name]
	if !ok {
		return fmt.Errorf("Could not find resource by name %s", name)
	}

	// Loop fields
	for _, f := range r.Fields {
		if err := e.handleField(object, f); err != nil {
			return err
		}
	}
	return nil
}

func (e *Editor) handleField(object interface{}, field Field) error {
	objectVal := reflect.ValueOf(object).Elem()
	fieldVal := objectVal.FieldByName(field.Name)
	kind := fieldVal.Kind()

	if kind == reflect.Slice {
		kind = reflect.TypeOf(fieldVal.Interface()).Elem().Kind()
		for {
			vi, blank, err := e.something(field, kind)
			if err != nil {
				return err
			}

			var value reflect.Value
			if blank {
				value = reflect.Zero(fieldVal.Type())
			}
			value = reflect.ValueOf(vi)
			fieldVal.Set(reflect.Append(fieldVal, value))

			return nil
		}
	} else {
		vi, blank, err := e.something(field, kind)
		if err != nil {
			return err
		}

		var value reflect.Value
		if blank {
			value = reflect.Zero(fieldVal.Type())
		} else {
			value = reflect.ValueOf(vi)
		}
		fieldVal.Set(value)
	}

	return nil
}

func (e *Editor) something(field Field, kind reflect.Kind) (interface{}, bool, error) {
	if kind == reflect.Struct {
		v, err := e.NewObject(field.Association)
		objectVal := reflect.ValueOf(v).Elem().Interface()
		return objectVal, false, err
	}

	return e.askAndConvert(field, kind)
}

func (e *Editor) askAndConvert(field Field, kind reflect.Kind) (interface{}, bool, error) {
	for {
		e.writef("%s (%s): ", field.Name, kind)
		text, err := e.getLine()
		if err != nil {
			return nil, false, err
		}

		if text == "" {
			return nil, true, nil
		}
		converted, err := e.convertValueType(text, kind)
		if err != nil {
			continue
		}
		return converted, false, nil
	}
}

func (e *Editor) convertValueType(value string, kind reflect.Kind) (interface{}, error) {
	switch kind {
	case reflect.String:
		return value, nil
	case reflect.Int:
		return strconv.Atoi(value)
	case reflect.Bool:
		return strconv.ParseBool(value)
	default:
		return nil, fmt.Errorf("Could not convert string to kind %s", kind)
	}
}

func (e *Editor) writef(format string, a ...interface{}) {
	fmt.Fprintf(e.writer, format, a...)
}

func (e *Editor) getLine() (string, error) {
	text, err := e.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}
