package whiteboard

import (
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
