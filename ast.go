package whiteboard

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/traefik/yaegi/interp"
)

var precedence = map[string]int{"+": 20, "-": 20, "*": 40, "/": 40, "%": 40, "^": 60}

type ExprAST interface {
	toStr() string
}

type NumberExprAST struct {
	Val float64
	Str string
}

type StrExprAST struct {
	Str string
}

type BinaryExprAST struct {
	Op string
	Lhs,
	Rhs ExprAST
}

type FunCallerExprAST struct {
	Name string
	Arg  []ExprAST
}

type SelectorExprAST struct {
	Name     string
	Selector Selector
}

func (n NumberExprAST) toStr() string {
	return fmt.Sprintf(
		"NumberExprAST:%s",
		n.Str,
	)
}

func (n StrExprAST) toStr() string {
	return fmt.Sprintf(
		"StrExprAST:%s",
		n.Str,
	)
}

func (b BinaryExprAST) toStr() string {
	return fmt.Sprintf(
		"BinaryExprAST: (%s %s %s)",
		b.Op,
		b.Lhs.toStr(),
		b.Rhs.toStr(),
	)
}

func (n FunCallerExprAST) toStr() string {
	return fmt.Sprintf(
		"FunCallerExprAST:%s",
		n.Name,
	)
}

func (s SelectorExprAST) toStr() string {
	return fmt.Sprintf(
		"SelectorExprAST:%s(%v)",
		s.Name, s.Selector,
	)
}

type AST struct {
	Tokens []*Token

	source    string
	currTok   *Token
	currIndex int
	depth     int

	Err error
}

func NewAST(toks []*Token, s string) *AST {
	a := &AST{
		Tokens: toks,
		source: s,
	}
	if a.Tokens == nil || len(a.Tokens) == 0 {
		a.Err = errors.New("empty token")
	} else {
		a.currIndex = 0
		a.currTok = a.Tokens[0]
	}
	return a
}

// Parser entry
func (a *AST) ParseExpression() ExprAST {
	a.depth++ // called depth
	lhs := a.parsePrimary()
	r := a.parseBinOpRHS(0, lhs)
	a.depth--
	if a.depth == 0 && a.currIndex != len(a.Tokens) && a.Err == nil {
		a.Err = errors.New(
			fmt.Sprintf("bad expression, reaching the end or missing the operator\n%s",
				ErrPos(a.source, a.currTok.Offset)))
	}
	return r
}

// Get the next Token
func (a *AST) getNextToken() *Token {
	a.currIndex++
	if a.currIndex < len(a.Tokens) {
		a.currTok = a.Tokens[a.currIndex]
		return a.currTok
	}
	return nil
}

// Get the operation priority
func (a *AST) getTokPrecedence() int {
	// fmt.Printf("getTokPrecedence-->%v\n", a.currTok.Tok)
	if p, ok := precedence[a.currTok.Tok]; ok {
		return p
	}
	return -1
}

// Parse the number and generate a NumberExprAST node
func (a *AST) parseNumber() NumberExprAST {
	f64, err := strconv.ParseFloat(a.currTok.Tok, 64)
	if err != nil {
		a.Err = errors.New(
			fmt.Sprintf("%v\nwant '(' or '0-9' but get '%s'\n%s",
				err.Error(),
				a.currTok.Tok,
				ErrPos(a.source, a.currTok.Offset)))
		return NumberExprAST{}
	}
	n := NumberExprAST{
		Val: f64,
		Str: a.currTok.Tok,
	}
	a.getNextToken()
	return n
}

func (a *AST) parseFunction() SelectorExprAST {

	name := a.currTok.Tok
	// fmt.Printf("parseFunction-->%v\n", name)

	selectorType := name
	a.getNextToken()
	// call custom function
	exprs := make([]ExprAST, 0)
	if a.currTok.Tok == "(" {

		a.getNextToken()

		if a.currTok.Tok == ")" {
			// function call without parameters
			// ignore the process of parameter resolution
		} else {
			exprs = append(exprs, a.ParseExpression())
			for a.currTok.Tok != ")" && a.getNextToken() != nil {
				if a.currTok.Type == COMMA {
					continue
				}
				exprs = append(exprs, a.ParseExpression())
			}
		}
	}
	// fmt.Printf("parms-->%v\n", exprs)

	var ifaceSlice []interface{}
	var err error

	for _, part := range exprs {
		fmt.Printf("%v-->%v-->%v\n", part, reflect.TypeOf(part).Kind(), reflect.TypeOf(part).Name())
		switch part.(type) {
		case StrExprAST:
			ifaceSlice = append(ifaceSlice, part.(StrExprAST).Str)
		case NumberExprAST:
			ps := part.(NumberExprAST).Str
			if strings.Contains(ps, ".") {
				pf, _ := strconv.ParseFloat(ps, 64)
				ifaceSlice = append(ifaceSlice, pf)
			} else {
				pf, _ := strconv.ParseInt(ps, 10, 64)
				ifaceSlice = append(ifaceSlice, pf)
			}
		case SelectorExprAST:
			ifaceSlice = append(ifaceSlice, part.(SelectorExprAST).Selector)
		}
	}

	// fmt.Printf("%s", ifaceSlice)
	s := SelectorExprAST{}
	switch selectorType {
	case "K":
		if len(ifaceSlice) == 1 {
			s.Name = selectorType
			s.Selector, _ = NewK(ifaceSlice[0])
		} else {
			a.Err = errors.New(
				fmt.Sprintf("Selector `%s` is out of limit\n%s",
					s.Name,
					ErrPos(a.source, a.currTok.Offset)))
		}

	case "S":
		s.Name = selectorType
		s.Selector, err = NewS(ifaceSlice[0:]...)
		if err != nil {
			a.Err = errors.New(
				fmt.Sprintf("Selector `%s` %s \n%s",
					s.Name,
					err.Error(),
					ErrPos(a.source, a.currTok.Offset)))
		}

	case "F":
		s.Name = selectorType
		// TODO: Complete extraction function string
		// expr, err := eval.Parse(ifaceSlice[0])
		// program, err := expr.Compile(parts[0], expr.Env(Env{}))

		i := interp.New(interp.Options{})
		v, err := i.Eval(ifaceSlice[0].(string))
		if err != nil {
			a.Err = errors.New(
				fmt.Sprintf("Selector `%s` %s \n%s",
					s.Name,
					err.Error(),
					ErrPos(a.source, a.currTok.Offset)))
		}
		fn := v.Interface().(func(interface{}, ...interface{}) interface{})
		s.Selector = NewF(fn, ifaceSlice[1:]...)
	case "ExpS":
		s.Name = selectorType
		s.Selector = NewExpressionSelector(ifaceSlice[0].(Selector), ifaceSlice[1].(Selector), ifaceSlice[2].(string))

	case "IF":
		s.Name = selectorType
		s.Selector = NewIf(ifaceSlice[0].(Selector), ifaceSlice[1].(Selector), ifaceSlice[2].(Selector))
	case "AL":
		s.Name = selectorType
		var selectors []Selector
		for _, s := range ifaceSlice {
			if sel, ok := s.(Selector); ok {
				selectors = append(selectors, sel)
			}
		}
		s.Selector = NewAlternation(selectors...)
	}

	// fmt.Printf("parseSelector-->%v\n", s)

	a.getNextToken()
	return s
}

func (a *AST) parseFunCallerOrConst() ExprAST {
	name := a.currTok.Tok
	a.getNextToken()
	// call func
	if a.currTok.Tok == "(" {
		f := FunCallerExprAST{}
		if _, ok := defFunc[name]; !ok {
			a.Err = errors.New(
				fmt.Sprintf("function `%s` is undefined\n%s",
					name,
					ErrPos(a.source, a.currTok.Offset)))
			return f
		}
		a.getNextToken()
		exprs := make([]ExprAST, 0)
		if a.currTok.Tok == ")" {
			// function call without parameters
			// ignore the process of parameter resolution
		} else {
			exprs = append(exprs, a.ParseExpression())
			for a.currTok.Tok != ")" && a.getNextToken() != nil {
				if a.currTok.Type == COMMA {
					continue
				}
				exprs = append(exprs, a.ParseExpression())
			}
		}
		def := defFunc[name]
		if def.argc >= 0 && len(exprs) != def.argc {
			a.Err = errors.New(
				fmt.Sprintf("wrong way calling function `%s`, parameters want %d but get %d\n%s",
					name,
					def.argc,
					len(exprs),
					ErrPos(a.source, a.currTok.Offset)))
		}
		a.getNextToken()
		f.Name = name
		f.Arg = exprs
		return f
	}
	// call const
	if v, ok := defConst[name]; ok {
		return NumberExprAST{
			Val: v,
			Str: strconv.FormatFloat(v, 'f', 0, 64),
		}
	} else {
		// a.Err = errors.New(
		// 	fmt.Sprintf("const `%s` is undefined\n%s",
		// 		name,
		// 		ErrPos(a.source, a.currTok.Offset)))
		// return NumberExprAST{}

		return StrExprAST{
			Str: name,
		}

	}
}

// Get a node and return ExprAST
// All possible types are processed here and the corresponding types are resolved
func (a *AST) parsePrimary() ExprAST {
	switch a.currTok.Type {
	case Identifier:
		return a.parseFunCallerOrConst()
	case Literal:
		return a.parseNumber()
	case Operator:
		if a.currTok.Tok == "(" {
			t := a.getNextToken()
			if t == nil {
				a.Err = errors.New(
					fmt.Sprintf("want '(' or '0-9' but get EOF\n%s",
						ErrPos(a.source, a.currTok.Offset)))
				return nil
			}
			e := a.ParseExpression()
			if e == nil {
				return nil
			}
			if a.currTok.Tok != ")" {
				a.Err = errors.New(
					fmt.Sprintf("want ')' but get %s\n%s",
						a.currTok.Tok,
						ErrPos(a.source, a.currTok.Offset)))
				return nil
			}
			a.getNextToken()
			return e
		} else if a.currTok.Tok == "-" {
			if a.getNextToken() == nil {
				a.Err = errors.New(
					fmt.Sprintf("want '0-9' but get '-'\n%s",
						ErrPos(a.source, a.currTok.Offset)))
				return nil
			}
			bin := BinaryExprAST{
				Op:  "-",
				Lhs: NumberExprAST{},
				Rhs: a.parsePrimary(),
			}
			return bin
		} else {
			return a.parseNumber()
		}
	case COMMA:
		a.Err = errors.New(
			fmt.Sprintf("want '(' or '0-9' but get %s\n%s",
				a.currTok.Tok,
				ErrPos(a.source, a.currTok.Offset)))
		return nil
	case FUCTION:
		return a.parseFunction()
	default:
		return nil
	}
}

// Loop to obtain the priority of the operator, recursing the higher priority into deeper nodes
// This is the most important algorithm for generating the correct AST structure, and it must be carefully read and understood
func (a *AST) parseBinOpRHS(execPrec int, lhs ExprAST) ExprAST {
	for {
		tokPrec := a.getTokPrecedence()
		if tokPrec < execPrec {
			return lhs
		}
		binOp := a.currTok.Tok
		if a.getNextToken() == nil {
			a.Err = errors.New(
				fmt.Sprintf("want '(' or '0-9' but get EOF\n%s",
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}
		rhs := a.parsePrimary()
		if rhs == nil {
			return nil
		}
		nextPrec := a.getTokPrecedence()
		if tokPrec < nextPrec {
			rhs = a.parseBinOpRHS(tokPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}
		lhs = BinaryExprAST{
			Op:  binOp,
			Lhs: lhs,
			Rhs: rhs,
		}
	}
}
