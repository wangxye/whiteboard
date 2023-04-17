package whiteboard

import (
	"fmt"
	"reflect"
)

type Transport struct {
	value   interface{}
	context map[interface{}]interface{}
}

func NewTransport(value interface{}, context map[interface{}]interface{}) *Transport {
	return &Transport{
		value:   value,
		context: context,
	}
}

func (t *Transport) FromSource(src interface{}) *Transport {
	switch s := src.(type) {
	case *Transport:
		return s
	default:
		return NewTransport(s, make(map[interface{}]interface{}))
	}
}

type BendingException struct {
	Message string
}

func (e *BendingException) Error() string {
	return e.Message
}

func Bend(mapping interface{}, source interface{}, args ...interface{}) (interface{}, error) {
	context := make(map[interface{}]interface{})
	if len(args) > 0 {
		context, _ = args[0].(map[interface{}]interface{})
	}
	transport := NewTransport(source, context)
	return _bend(mapping, transport)
}

func _bend(mapping interface{}, transport *Transport) (interface{}, error) {
	t := reflect.TypeOf(mapping)
	// fmt.Println(t)
	fmt.Printf("%v-->%v\n", t, mapping)
	switch m := mapping.(type) {
	case []interface{}:
		result := make([]interface{}, len(m))
		for i, item := range m {
			val, err := _bend(item, transport)
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil
	case map[string]interface{}:
		result := make(map[interface{}]interface{})
		for k, v := range m {
			val, err := _bend(v, transport)
			if err != nil {
				return nil, &BendingException{
					Message: fmt.Sprintf("Error for key %v: %v", k, err.Error()),
				}
			}
			result[k] = val
		}
		fmt.Printf("_bend: %v\n", result)
		return result, nil
	case string:
		val, err := bendExpression(m, transport)
		fmt.Printf("_bend string: %v\n", val)
		return val, err

	default:
		return mapping, nil
	}
}

func bendExpression(mapping interface{}, transport *Transport) (interface{}, error) {
	exp := mapping.(string)
	toks, err := Parse(exp)

	if err != nil {
		return nil, &BendingException{
			Message: fmt.Sprintf("Error for lexical analysis: mapping: %v, error: %v", mapping, err.Error()),
		}
	}
	// []token -> AST Tree
	ast := NewAST(toks, exp)
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return nil, &BendingException{
			Message: fmt.Sprintf("Error for NewAst: mapping: %v, error: %v", mapping, ast.Err.Error()),
		}
	}

	// AST builder
	ar := ast.ParseExpression()
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return nil, &BendingException{
			Message: fmt.Sprintf("Error for AST builder: mapping: %v, error: %v", mapping, ast.Err.Error()),
		}
	}

	fmt.Printf("ExprAST: %+v\n", ar)

	// AST traversal -> result
	r, err := ExprASTResultWithContext(ar, transport.value)

	if r == nil && len(transport.context) != 0 {
		r, err = ExprASTResultWithContext(ar, transport.context)
	}

	fmt.Println("progressing ...\t", r)
	fmt.Printf("%s = %v\n", exp, r)

	return r, err

}
