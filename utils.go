package whiteboard

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func IsMapKeyTypeEqual(m reflect.Value, key reflect.Value) bool {
	// check whether m is a map type
	if m.Kind() != reflect.Map {
		return false
	}

	// gets the key type of the map
	mapKeyType := m.Type().Key()

	// check whether the type of the key is the same as that of the map key
	return key.Type().AssignableTo(mapKeyType)
}

func IsValidMatch(v reflect.Value, key reflect.Value) bool {
	// check whether v is a map type

	if v.Kind() == reflect.Map {
		// gets the key type of the map
		mapKeyType := v.Type().Key()

		// check whether the type of the key is the same as that of the map key
		return key.Type().AssignableTo(mapKeyType)
	}
	// check whether v is an array
	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {

		if key.Kind() != reflect.Int && key.Kind() != reflect.Int64 && key.Kind() != reflect.Int32 {
			fmt.Printf("%v -> %v return %t \n", v.Kind(), key.Kind(), key.Kind() == reflect.Int64)
			return false
		}

		keyInt := int(key.Int())

		if key.IsValid() && keyInt >= 0 && keyInt < v.Len() {
			fmt.Printf("%v -> %v return %t  %v\n", v.Kind(), key.Kind(), key.Type().Kind() == reflect.Int, keyInt)
			return true
		}

	}
	fmt.Printf("%v/ %v -> %v /%v \n", key.Type(), key, v.Type(), v)
	return false
}

// Top level function
// Analytical expression and execution
// err is not nil if an error occurs (including arithmetic runtime errors)
func ParseAndExec(s string) (r float64, err error) {
	toks, err := Parse(s)
	if err != nil {
		return 0, err
	}
	ast := NewAST(toks, s)
	if ast.Err != nil {
		return 0, ast.Err
	}
	ar := ast.ParseExpression()
	if ast.Err != nil {
		return 0, ast.Err
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	return ExprASTResult(ar).(float64), err
}

func ErrPos(s string, pos int) string {
	r := strings.Repeat("-", len(s)) + "\n"
	s += "\n"
	for i := 0; i < pos; i++ {
		s += " "
	}
	s += "^\n"
	return r + s + r
}

// the integer power of a number
func Pow(x float64, n float64) float64 {
	return math.Pow(x, n)
}

func expr2Radian(expr ExprAST) float64 {
	r := ExprASTResult(expr).(float64)
	if TrigonometricMode == AngleMode {
		r = r / 180 * math.Pi
	}
	return r
}

// Float64ToStr float64 -> string
func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// RegFunction is Top level function
// register a new function to use in expressions
// name: be register function name. the same function name only needs to be registered once.
// argc: this is a number of parameter signatures. should be -1, 0, or a positive integer
//
//	-1 variable-length argument; >=0 fixed numbers argument
//
// fun:  function handler
func RegFunction(name string, argc int, fun func(...ExprAST) float64) error {
	if len(name) == 0 {
		return errors.New("RegFunction name is not empty")
	}
	if argc < -1 {
		return errors.New("RegFunction argc should be -1, 0, or a positive integer")
	}
	if _, ok := defFunc[name]; ok {
		return errors.New("RegFunction name is already exist")
	}
	defFunc[name] = defS{argc, fun}
	return nil
}

// ExprASTResult is a Top level function
// AST traversal
// if an arithmetic runtime error occurs, a panic exception is thrown
func ExprASTResult(expr ExprAST) interface{} {
	var l, r interface{}

	//TODO: handle error and return a uniform result type
	fmt.Printf("ExprASTResult-->%v\n", expr)

	switch expr.(type) {
	case BinaryExprAST:
		ast := expr.(BinaryExprAST)
		l = ExprASTResult(ast.Lhs)
		r = ExprASTResult(ast.Rhs)
		switch ast.Op {
		case "+":
			// strconv.Atoi(l)
			switch v := l.(type) {
			case int64, int, int16:
				il := l.(int64)
				ir := r.(int64)
				return il + ir
			case float64:
				fl := l.(float64)
				fr := r.(float64)
				return fl + fr
			case string:
				sl := l.(string)
				sr := r.(string)
				return sl + sr
			default:
				panic(fmt.Sprintf("unsupported type %T in addition operation", v))
			}

		case "-":
			fl := l.(float64)
			fr := r.(float64)
			return fl - fr
		case "*":
			fl := l.(float64)
			fr := r.(float64)
			return fl * fr
		case "/":
			if r.(float64) == 0 {
				panic(errors.New(
					fmt.Sprintf("violation of arithmetic specification: a division by zero in ExprASTResult: [%g/%g]",
						l,
						r)))
			}
			fl := l.(float64)
			fr := r.(float64)
			return fl / fr
		case "%":
			if r.(float64) == 0 {
				panic(errors.New(
					fmt.Sprintf("violation of arithmetic specification: a division by zero in ExprASTResult: [%g%%%g]",
						l,
						r)))
			}
			il := int(l.(float64))
			ir := int(r.(float64))
			return il % ir
		case "^":
			fl := l.(float64)
			fr := r.(float64)
			res := Pow(fl, fr)
			return res
		default:
			panic(fmt.Sprintf("unsupported operator %s", ast.Op))
		}
	case NumberExprAST:
		return expr.(NumberExprAST).Val
	case FunCallerExprAST:
		f := expr.(FunCallerExprAST)
		def := defFunc[f.Name]
		return def.fun(f.Arg...)
	case SelectorExprAST:
		sea := expr.(SelectorExprAST)
		// r, _ := sea.Selector.Execute(nil)
		// fmt.Printf("%v-->%v\n", r, reflect.TypeOf(sea.Selector))
		// return r
		var r interface{}
		switch v := sea.Selector.(type) {
		case *K:
			r, _ = sea.Selector.Execute(nil)
		default:
			panic(fmt.Sprintf("unsupported type %T in addition operation", v))
		}
		fmt.Printf("%v-->%v\n", r, reflect.TypeOf(sea.Selector))

		return r
	}

	return nil
}

func ExprASTResultWithContext(expr ExprAST, context interface{}) (interface{}, error) {
	var l, r interface{}
	// var err error
	//TODO: handle error and return a uniform result type
	fmt.Printf("ExprASTResult-->%v\n", expr)

	switch expr.(type) {
	case BinaryExprAST:
		ast := expr.(BinaryExprAST)
		l = ExprASTResult(ast.Lhs)
		r = ExprASTResult(ast.Rhs)
		switch ast.Op {
		case "+":
			// strconv.Atoi(l)
			switch v := l.(type) {
			case int64, int, int16:
				il := l.(int64)
				ir := r.(int64)
				return il + ir, nil
			case float64, float32:
				fl := l.(float64)
				fr := r.(float64)
				return fl + fr, nil
			case string:
				sl := l.(string)
				sr := r.(string)
				return sl + sr, nil
			default:
				panic(fmt.Sprintf("unsupported type %T in addition operation", v))
			}

		case "-":
			fl := l.(float64)
			fr := r.(float64)
			return fl - fr, nil
		case "*":
			fl := l.(float64)
			fr := r.(float64)
			return fl * fr, nil
		case "/":
			if r.(float64) == 0 {
				errMsg := fmt.Sprintf("violation of arithmetic specification: a division by zero in ExprASTResult: [%g/%g]", l, r)
				return nil, errors.New(errMsg)
			}
			fl := l.(float64)
			fr := r.(float64)
			return fl / fr, nil
		case "%":
			if r.(float64) == 0 {
				errMsg := fmt.Sprintf("violation of arithmetic specification: a division by zero in ExprASTResult: [%g/%g]", l, r)
				return nil, errors.New(errMsg)
			}
			il := int(l.(float64))
			ir := int(r.(float64))
			return il % ir, nil
		case "^":
			fl := l.(float64)
			fr := r.(float64)
			res := Pow(fl, fr)
			return res, nil
		default:
			panic(fmt.Sprintf("unsupported operator %s", ast.Op))
		}
	case NumberExprAST:
		return expr.(NumberExprAST).Val, nil
	case FunCallerExprAST:
		f := expr.(FunCallerExprAST)
		def := defFunc[f.Name]
		return def.fun(f.Arg...), nil
	case SelectorExprAST:
		sea := expr.(SelectorExprAST)
		// var r interface{}
		r, err := sea.Selector.Execute(context)
		return r, err
	}

	return nil, nil
}
