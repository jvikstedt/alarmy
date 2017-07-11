package transform

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)
var errQuit = errors.New("quit")

type Field struct {
	ID     int
	Name   string
	Value  interface{}
	Kind   reflect.Kind
	Locked bool
}

func Fields(i interface{}) []Field {
	fields := []Field{}

	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for id := 0; id < v.NumField(); id++ {
		field := Field{}

		field.ID = id
		field.Name = v.Type().Field(id).Name
		field.Value = v.Field(id).Interface()
		field.Kind = v.Field(id).Kind()
		field.Locked = Locked(i, id)

		fields = append(fields, field)
	}

	return fields
}

func printFields(fields []Field) {
	for _, v := range fields {
		fmt.Printf("%d: %s = %v\n", v.ID, v.Name, v.Value)
	}
}

func getLine() (string, error) {
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func askFieldByID(fields []Field) (Field, error) {
	for {
		fmt.Printf("\nType field id to edit: ")

		line, err := getLine()
		if err != nil {
			return Field{}, err
		}

		switch line {
		case "q":
			return Field{}, errQuit
		default:
			id, err := strconv.Atoi(line)
			if err != nil {
				fmt.Println(err)
				break
			}

			field, err := fieldByID(fields, id)
			if err != nil {
				fmt.Println(err)
				break
			}

			if field.Locked {
				fmt.Println("Can't edit locked field")
				break
			}

			return field, nil
		}
	}
}

func askValue() (string, error) {
	fmt.Printf("Type value: ")
	return getLine()
}

func fieldByID(fields []Field, id int) (Field, error) {
	for _, v := range fields {
		if v.ID == id {
			return v, nil
		}
	}
	return Field{}, fmt.Errorf("Field with id of %d not found", id)
}

func Modify(i interface{}) error {
	for {
		fields := Fields(i)
		printFields(fields)

		field, err := askFieldByID(fields)
		if err != nil {
			if err != errQuit {
				return err
			}

			continue
		}

		value, err := askValue()
		if err != nil {
			return err
		}

		realVal, err := convertVal(field, value)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = SetValue(i, field.ID, realVal)
		if err != nil {
			fmt.Println(err)
			continue
		}

		quit, err := askQuit()
		if err != nil {
			return err
		}
		fmt.Println()

		if quit {
			break
		}
	}

	return nil
}

func convertVal(field Field, value string) (interface{}, error) {
	switch field.Kind {
	case reflect.String:
		return value, nil
	case reflect.Int:
		return strconv.Atoi(value)
	default:
		return nil, fmt.Errorf("could not convert value %s to type %v", value, field.Kind)
	}
}

func askQuit() (bool, error) {
	for {
		fmt.Print("\nWant to quit editing? (y): ")
		line, err := getLine()

		if line == "y" {
			return true, err
		} else {
			return false, err
		}
	}
}

func SetValue(i interface{}, id int, value interface{}) error {

	valueKind := reflect.TypeOf(value).Kind()

	targetStruct := reflect.ValueOf(i)
	if targetStruct.Kind() == reflect.Ptr {
		targetStruct = targetStruct.Elem()
	} else {
		return fmt.Errorf("expected pointer")
	}

	if Locked(targetStruct.Interface(), id) {
		return fmt.Errorf("can't set value of locked field")
	}
	fieldValue := reflect.ValueOf(value)

	targetKind := targetStruct.Field(id).Kind()

	if valueKind != targetKind {
		return fmt.Errorf("target and source value does not match %s != %s", targetKind, valueKind)
	}

	targetStruct.Field(id).Set(fieldValue)

	return nil
}

func Locked(i interface{}, id int) bool {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	structTag := v.Type().Field(id).Tag

	tags, _ := structTag.Lookup("transform")
	return strings.Contains(tags, "lock")
}
