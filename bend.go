package whiteboard

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
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

var (
	logger *logrus.Entry
)

type BendingException struct {
	Message string
}

func (e *BendingException) Error() string {
	return e.Message
}

func Bend(mapping interface{}, source interface{}, args ...interface{}) (interface{}, error) {
	// check whether mapping and source are empty
	if mapping == nil || source == nil {
		return nil, errors.New("mapping or source is empty")
	}

	context := make(map[interface{}]interface{})
	// logger.Info("Bending source with mapping", source, mapping)
	if len(args) > 0 {
		context, _ = args[0].(map[interface{}]interface{})
	}
	transport := NewTransport(source, context)
	return _bend(mapping, transport)
}

func _bend(mapping interface{}, transport *Transport) (interface{}, error) {

	t := reflect.TypeOf(mapping)
	fmt.Println(t.Kind())
	mValue := reflect.ValueOf(mapping)

	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		result := make([]interface{}, mValue.Len())
		for i := 0; i < mValue.Len(); i++ {
			item := mValue.Index(i)
			val, err := _bend(item, transport)
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil
	case reflect.Map:
		// result := make(map[string]interface{})
		result := reflect.New(t).Elem()
		result.Set(reflect.MakeMap(t))
		keys := mValue.MapKeys()
		for _, key := range keys {

			val, err := _bend(mValue.MapIndex(key).Interface(), transport)
			if err != nil {
				return nil, &BendingException{
					Message: fmt.Sprintf("Error for key %v: %v", key, err.Error()),
				}
			}
			// result[key.Interface().(string)] = val
			valValue := reflect.ValueOf(val)
			result.SetMapIndex(key, valValue)
		}
		return result.Interface(), nil
	case reflect.String:
		val, err := bendExpression(mValue.Interface(), transport)
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
