package whiteboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type selector interface {
	Execute(source interface{}) interface{}
}

type K struct {
	Value interface{}
}

func NewK(value interface{}) *K {
	return &K{Value: value}
}

func (k *K) Execute(source interface{}) interface{} {
	return k.Value
}

type S struct {
	Path []interface{}
}

func NewS(path ...interface{}) (*S, error) {
	if len(path) == 0 {
		return nil, errors.New("No path given")
	}
	return &S{Path: path}, nil
}

func (s *S) Execute(source interface{}) (interface{}, error) {

	v := reflect.ValueOf(source)
	for _, key := range s.Path {
		if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
			if v.IsNil() {
				return reflect.Value{}, fmt.Errorf("nil pointer encountered in path")
			}
			v = v.Elem()
		}
		switch v.Kind() {
		case reflect.Struct:
			field := v.FieldByName(key.(string))
			if !field.IsValid() {
				return reflect.Value{}, fmt.Errorf("no such field %s", key)
			}
			v = field
		case reflect.Slice, reflect.Array:
			index := key.(int)
			if index < 0 || index >= v.Len() {
				return reflect.Value{}, fmt.Errorf("index out of range: %d", index)
			}
			v = v.Index(index)
			// break
		case reflect.Map:
			key := reflect.ValueOf(key)
			elem := v.MapIndex(key)
			if !elem.IsValid() {
				return reflect.Value{}, fmt.Errorf("no such key %s", key)
			}
			v = elem
		default:
			return reflect.Value{}, fmt.Errorf("cannot access path element %s of non-composite type %s", key, v.Kind())
		}
	}
	return v.Interface(), nil
}

type F struct {
	Func func(interface{}, ...interface{}) interface{}
	Args []interface{}
	// Kwargs map[interface{}]interface{}
}

func NewF(f func(interface{}, ...interface{}) interface{}, args ...interface{}) *F {
	return &F{Func: f, Args: args}
}

func (f *F) Execute(value interface{}) interface{} {
	args := f.Args
	if len(args) == 0 { // add this line to check if args is empty
		args = make([]interface{}, 1)
	}
	args = append([]interface{}{value}, args...)
	if len(args) == 1 {
		return f.Func(args[0])
	}
	return f.Func(args[0], args[1:]...)
}
