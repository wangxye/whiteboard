package whiteboard

import (
	"errors"
	"fmt"
)

type If struct {
	condition Selector
	whenTrue  Selector
	whenFalse Selector
}

func (i *If) Execute(val interface{}) (interface{}, error) {

	condVal, err := i.condition.Execute(val)
	if err != nil {
		return nil, err
	}
	if condVal.(bool) {
		return i.whenTrue.Execute(val)
	} else {
		return i.whenFalse.Execute(val)
	}
}

func NewIf(condition Selector, whenTrue Selector, whenFalse Selector) *If {
	return &If{condition: condition, whenTrue: whenTrue, whenFalse: whenFalse}
}

type Alternation struct {
	selectors []Selector
}

func NewAlternation(s ...Selector) *Alternation {
	return &Alternation{s}
}

func (a *Alternation) Execute(source interface{}) (interface{}, error) {
	var exc error
	for _, selector := range a.selectors {
		result, err := selector.Execute(source)
		// && !errors.Is(err, NotFoundError)
		fmt.Printf("%v -> %v -> %v \n", source, result, err)
		if err != nil {
			exc = err
		} else {
			return result, nil
		}
	}
	return nil, exc
}

type Switch struct {
	keySelctor      Selector
	cases           map[interface{}]Selector
	defaultSelector Selector
}

func (s *Switch) Execute(source map[interface{}]interface{}) (interface{}, error) {

	key, err := s.keySelctor.Execute(source)
	if err != nil {
		return nil, err
	}

	var benderFn Selector
	var ok bool
	if benderFn, ok = s.cases[key]; !ok {
		if s.defaultSelector != nil {
			benderFn = s.defaultSelector
		} else {
			return nil, errors.New("key not found in case container")
		}
	}

	val, err := benderFn.Execute(source)
	if err != nil {
		return nil, err
	}

	return val, nil
}
