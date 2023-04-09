package whiteboard

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_FunctionStr_Extract(t *testing.T) {
	str := `func(v interface{}, args ...interface{}) interface{} {
		a := args[0].(string)
		b := args[1].(int)
		return fmt.Sprintf("%s%d", a, b+v.(int))
	}, "suffix", 100`

	// Split string by comma
	// parts := strings.Split(str, ",")

	// Find function string
	// var function string
	// for _, part := range parts {
	// 	if strings.Contains(part, "func") {
	// 		function = strings.TrimSpace(part)
	// 		break
	// 	}
	// }
	pattern := `func\([^)]*\)[^{]*\{(?:[^{}]|(?R))*\}`
	re := regexp.MustCompile(pattern)
	function := re.FindString(str)

	expect := `func(v interface{}, args ...interface{}) interface{} {
		a := args[0].(string)
		b := args[1].(int)
		return fmt.Sprintf("%s%d", a, b+v.(int))
	}`
	// Print result
	// fmt.Println("Function:", function)
	if function != expect {
		t.Errorf("Expected Fuction, but got %s", function)
	}
}

func Test_Fuction_AST_Binary(t *testing.T) {
	exp := "92 + 5 + 5 * 27 - (92 - 12) / 4 + 26"
	// input text -> []token
	toks, err := Parse(exp)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	// []token -> AST Tree
	ast := NewAST(toks, exp)
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	// AST builder
	ar := ast.ParseExpression()
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	fmt.Printf("ExprAST: %+v\n", ar)
	// 加入下面的代码
	// AST traversal -> result
	r := ExprASTResult(ar)
	fmt.Println("progressing ...\t", r)
	fmt.Printf("%s = %v\n", exp, r)

	expect := 238.0
	if r != expect {
		t.Errorf("Expected Fuction, but got %f", r)
	}
}

func Test_Fuction_AST_Selector(t *testing.T) {
	exp := "K(1) + K(2)"
	// input text -> []token
	toks, err := Parse(exp)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	// []token -> AST Tree
	ast := NewAST(toks, exp)
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	// AST builder
	ar := ast.ParseExpression()
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	fmt.Printf("ExprAST: %+v\n", ar)
	// 加入下面的代码
	// AST traversal -> result
	r := ExprASTResult(ar)
	fmt.Println("progressing ...\t", r)
	fmt.Printf("%s = %v\n", exp, r)

	expect := 3
	val := int(r.(float64))
	if val != expect {
		t.Errorf("Expected Fuction, but got %f", r)
	}
}

func Test_Fuction_AST_Selector_str(t *testing.T) {
	exp := "K(\"1\") + K(\"2\")"
	// input text -> []token
	toks, err := Parse(exp)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		return
	}
	// []token -> AST Tree
	ast := NewAST(toks, exp)
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	// AST builder
	ar := ast.ParseExpression()
	if ast.Err != nil {
		fmt.Println("ERROR: " + ast.Err.Error())
		return
	}
	fmt.Printf("ExprAST: %+v\n", ar)
	// 加入下面的代码
	// AST traversal -> result
	r := ExprASTResult(ar)
	fmt.Println("progressing ...\t", r)
	fmt.Printf("%s = %v\n", exp, r)

	expect := "12"
	if r != expect {
		t.Errorf("Expected Fuction, but got %f", r)
	}
}
