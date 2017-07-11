package transform_test

import (
	"reflect"
	"testing"

	"github.com/jvikstedt/alarmy/transform"
	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	obj := struct {
		ID   int
		Name string
	}{ID: 1, Name: "Golang"}

	fields := transform.Fields(&obj)

	expectedResult := []transform.Field{
		transform.Field{ID: 0, Name: "ID", Value: 1, Kind: reflect.Int, Locked: false},
		transform.Field{ID: 1, Name: "Name", Value: "Golang", Kind: reflect.String, Locked: false},
	}

	assert.Equal(t, expectedResult, fields)
}

func TestSetValue(t *testing.T) {
	obj := struct {
		ID   int
		Name string
	}{}

	transform.SetValue(&obj, 0, 15)
	assert.Equal(t, 15, obj.ID)

	transform.SetValue(&obj, 1, "something")
	assert.Equal(t, "something", obj.Name)

	err := transform.SetValue(&obj, 0, "something")
	assert.Error(t, err)
	assert.Equal(t, 15, obj.ID)

	lockObj := struct {
		ID int `transform:"lock"`
	}{}

	err = transform.SetValue(&lockObj, 0, 15)
	assert.Error(t, err)
	assert.Equal(t, 0, lockObj.ID)
}

func TestLocked(t *testing.T) {
	obj := struct {
		ID int
	}{}

	locked := transform.Locked(obj, 0)
	assert.False(t, locked, "expecting false when no tag value lock")

	objLock := struct {
		ID int `transform:"lock"`
	}{}

	locked = transform.Locked(objLock, 0)
	assert.True(t, locked, "expecting return true when has tag value lock")
}
