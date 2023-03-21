package whiteboard

type If struct {
	condition func(interface{}) bool
	whenTrue  selector
	whenFalse selector
}

func (i *If) Execute(val interface{}) (interface{}, error) {
	if i.condition(val) {
		return i.whenTrue.Execute(val)
	} else {
		return i.whenFalse.Execute(val)
	}
}

func NewIf(condition func(interface{}) bool, whenTrue selector, whenFalse selector) *If {
	return &If{condition: condition, whenTrue: whenTrue, whenFalse: whenFalse}
}

type Alternation struct {
	selectors []selector
}

func NewAlternation(s ...selector) *Alternation {
	return &Alternation{s}
}

func (a *Alternation) Execute(source interface{}) (interface{}, error) {
	var exc error
	for _, selector := range a.selectors {
		result, err := selector.Execute(source)
		// && !errors.Is(err, NotFoundError)
		if err != nil {
			exc = err
		} else {
			return result, nil
		}
	}
	return nil, exc
}
