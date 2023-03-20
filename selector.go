package whiteboard

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (s *S) Execute(source interface{}) interface{} {
	for _, key := range s.Path {
		if m, ok := source.(map[interface{}]interface{}); ok {
			source = m[key]
		} else if s, ok := source.([]map[interface{}]interface{}); ok && key == 0 {
			// Special case for the first key when the source is a slice of maps
			source = s[0]
		} else {
			// Handle unexpected value types, e.g. by converting to JSON
			jsonStr, _ := json.Marshal(source)
			panic(fmt.Sprintf("Cannot convert %v to map: %s", source, jsonStr))
		}
	}
	return source

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
