package whiteboard

import "errors"

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
		source = source.(map[interface{}]interface{})[key]
	}
	return source
}

/**

type K struct {
	val interface{}
}

func (k K) Execute(source interface{}) interface{} {
	return k.val
}

type S struct {
	path interface{}
}

func (s *S) Execute(source interface{}) interface{} {
	path, ok := s.path.([]interface{})
	if !ok {
		return nil
	}

	for _, key := range path {
		if m, ok := source.(map[string]interface{}); ok {
			value, exists := m[key.(string)]
			if exists {
				source = value
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	return source
}
**/

type F struct {
	Func   func(interface{}, ...interface{}) interface{}
	Args   []interface{}
	Kwargs map[interface{}]interface{}
}

func NewF(f func(interface{}, ...interface{}) interface{}, args ...interface{}) *F {
	kwargs := map[interface{}]interface{}{}
	for i := 0; i < len(args); i += 2 {
		kwargs[args[i]] = args[i+1]
	}
	return &F{Func: f, Args: args, Kwargs: kwargs}
}

func (f *F) Execute(value interface{}) interface{} {
	var args []interface{}
	for _, arg := range f.Args {
		if _, ok := f.Kwargs[arg]; !ok {
			args = append(args, arg)
		}
	}
	args = append([]interface{}{value}, args...)
	return f.Func(args)
}

//
