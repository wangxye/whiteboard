package whiteboard

type If struct {
	condition func(interface{}) bool
	whenTrue  selector
	whenFalse selector
}

func (i *If) Execute(val interface{}) interface{} {
	if i.condition(val) {
		return i.whenTrue.Execute(val)
	} else {
		return i.whenFalse.Execute(val)
	}
}

func NewIf(condition func(interface{}) bool, whenTrue selector, whenFalse selector) *If {
	return &If{condition: condition, whenTrue: whenTrue, whenFalse: whenFalse}
}
