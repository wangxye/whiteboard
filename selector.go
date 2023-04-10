package whiteboard

import (
	"errors"
	"fmt"
	"reflect"
)

type selector interface {
	Execute(source interface{}) (interface{}, error)
}

type K struct {
	Value interface{}
}

func NewK(value interface{}) (*K, error) {
	return &K{Value: value}, nil
}

func (k *K) Execute(source interface{}) (interface{}, error) {
	return k.Value, nil
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
	if source == nil {
		return nil, fmt.Errorf("KeyError:invalid reflect.Value")
	}
	v := reflect.ValueOf(source)

	for _, key := range s.Path {
		field, err := s.findFieldByKind(v, key)
		if err != nil {
			return nil, err
		}
		v = field
	}

	return v.Interface(), nil
}

func (s *S) findFieldByKind(v reflect.Value, key interface{}) (reflect.Value, error) {

	k := reflect.ValueOf(key)
	fmt.Printf("%v -- > %v\n", key, v)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}, fmt.Errorf("nil encountered in path")
		}
		v = v.Elem()
	}

	if !IsValidMatch(v, k) {
		return reflect.Value{}, fmt.Errorf("type inconsistency %s- > %s", k.Kind(), v.Kind())
	}

	switch v.Kind() {
	case reflect.Struct:
		field := v.FieldByName(key.(string))
		if !field.IsValid() {
			return reflect.Value{}, fmt.Errorf("no such field %s", key)
		}
		v = field
		// TODO: A field is a structure whose internal fields are recursively accessed
		// if v.Type().Kind() == reflect.Struct {
		// 	var err error
		// 	v, err = s.findFieldByKind(v, key)
		// 	if err != nil {
		// 		return reflect.Value{}, err
		// 	}
		// }
	case reflect.Slice, reflect.Array:
		// TODO:	interface {} is string, not int
		// index := key.(int)
		index := int(k.Int())
		if index < 0 || index >= v.Len() {
			return reflect.Value{}, fmt.Errorf("index out of range: %d", index)
		}
		v = v.Index(index)

	case reflect.Map:
		//value of type int is not assignable to type string
		// key := reflect.ValueOf(key)
		// fmt.Print(key.Kind())
		elem := v.MapIndex(k)
		if !elem.IsValid() {
			return reflect.Value{}, fmt.Errorf("no such key %s", key)
		}
		v = elem
	default:
		return reflect.Value{}, fmt.Errorf("cannot access path element %s of non-composite type %s", key, v.Kind())
	}
	return v, nil
}

type F struct {
	Func func(interface{}, ...interface{}) interface{}
	Args []interface{}
	// Kwargs map[interface{}]interface{}
}

func NewF(f func(interface{}, ...interface{}) interface{}, args ...interface{}) *F {
	return &F{Func: f, Args: args}
}

func (f *F) Execute(value interface{}) (interface{}, error) {
	args := f.Args
	if len(args) == 0 { // add this line to check if args is empty
		args = make([]interface{}, 1)
	}
	args = append([]interface{}{value}, args...)
	if len(args) == 1 {
		return f.Func(args[0]), nil
	}
	return f.Func(args[0], args[1:]...), nil
}
