package edit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Kind uint

const (
	String Kind = iota
	Int
	Bool
)

type FieldEditor func(interface{}, Field) error

type Field struct {
	Name string
	Kind
	Default interface{}
	FieldEditor
}

func KindTranslation(kind Kind) string {
	switch kind {
	case String:
		return "String"
	case Int:
		return "Integer"
	case Bool:
		return "Boolean"
	default:
		return "Unknown Type"
	}
}

func DefaultValueFor(field Field) (interface{}, error) {
	if field.Default != nil {
		return field.Default, nil
	}
	switch field.Kind {
	case String:
		return "", nil
	case Int:
		return 0, nil
	case Bool:
		return false, nil
	default:
		return nil, fmt.Errorf("Could not find default value for %v", field.Kind)
	}
}

func ConvertValueType(value string, field Field) (interface{}, error) {
	if value == "" {
		return DefaultValueFor(field)
	}

	switch field.Kind {
	case String:
		return value, nil
	case Int:
		return strconv.Atoi(value)
	case Bool:
		return strconv.ParseBool(value)
	default:
		return nil, fmt.Errorf("Could not convert string to kind %v", field.Kind)
	}
}

func SetObjectField(object interface{}, field Field, value string) error {
	objectVal := reflect.ValueOf(object).Elem()
	fieldVal := objectVal.FieldByName(field.Name)

	converted, err := ConvertValueType(value, field)
	if err != nil {
		return err
	}

	convertedVal := reflect.ValueOf(converted)

	fieldVal.Set(convertedVal)

	return nil
}

func EditObjectField(object interface{}, field Field) error {
	fmt.Printf("%s (%s): ", field.Name, KindTranslation(field.Kind))
	for {
		text, err := getLine()
		if err != nil {
			return err
		}

		if err := SetObjectField(object, field, text); err != nil {
			fmt.Println(err)
			continue
		}

		break
	}
	return nil
}

func Edit(object interface{}, fields []Field) error {
	for _, f := range fields {
		if f.FieldEditor != nil {
			if err := f.FieldEditor(object, f); err != nil {
				return err
			}
			continue
		}
		if err := EditObjectField(object, f); err != nil {
			return err
		}
	}

	return nil
}

func ObjectPrettyFormat(object interface{}) (string, error) {
	data, err := json.MarshalIndent(object, "", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

var reader = bufio.NewReader(os.Stdin)

func getLine() (string, error) {
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}
