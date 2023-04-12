package whiteboard

import (
	"fmt"
	"testing"
)

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

	// AST traversal -> result
	r := ExprASTResult(ar)
	fmt.Println("progressing ...\t", r)
	fmt.Printf("%s = %v\n", exp, r)

	expect := 3
	val := int(r.(int64))
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

	// AST traversal -> result
	r := ExprASTResult(ar)
	fmt.Println("progressing ...\t", r)
	fmt.Printf("%s = %v\n", exp, r)

	expect := "12"
	if r != expect {
		t.Errorf("Expected Fuction, but got %f", r)
	}
}
