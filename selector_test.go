package whiteboard

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestK_Execute(t *testing.T) {
	k, _ := NewK("Hello World")
	result, _ := k.Execute(nil)
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', but got %v", result)
	}
}

func TestS_Execute(t *testing.T) {
	data := map[interface{}]interface{}{
		"a": []map[interface{}]interface{}{
			{"b": 42},
			{"c": 32},
		},
	}

	s, err := NewS("a", 0, "b")
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}

	result, _ := s.Execute(data)
	if result != 42 {
		t.Errorf("Expected 42, but got %v", result)
	}
}

func TestS_Execute_single(t *testing.T) {
	data := map[interface{}]interface{}{
		"a": []map[interface{}]interface{}{
			{"b": 42},
			{"c": 32},
		},
	}

	s, err := NewS("a", 1, "c")
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}

	result, _ := s.Execute(data)
	if result != 32 {
		t.Errorf("Expected 32, but got %v", result)
	}
}

func TestS_Execute_str(t *testing.T) {
	// test case 2
	s2, err := NewS("a", "b", "c")

	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}

	source2 := map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": "hello"}}}
	res2, _ := s2.Execute(source2)
	if res2 != "hello" {
		t.Errorf("Test case 2 failed: expected %v but got %v", "hello", res2)
	}
}

func TestS_Execute_int(t *testing.T) {
	s1, err := NewS("a", 0, "b")
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}
	source1 := map[string]interface{}{"a": []interface{}{map[string]interface{}{"b": 42}}}
	res1, err := s1.Execute(source1)
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}
	if res1 != 42 {
		t.Errorf("Test case 1 failed: expected %v but got %v", 42, res1)
	}
}

func mySortFunc(args ...interface{}) interface{} {
	value := args[0]
	key := args[1].(string)
	list := value.([]map[string]interface{})
	sort.Slice(list, func(i, j int) bool {
		return list[i][key].(int) < list[j][key].(int)
	})
	return list
}

func TestF_Execute(t *testing.T) {
	data := []map[string]interface{}{{"id": 3}, {"id": 1}}

	f := NewF(func(value interface{}, args ...interface{}) interface{} {
		return mySortFunc(append([]interface{}{value}, args...)...)
	}, "id")

	result, _ := f.Execute(data)
	expected := []map[string]interface{}{{"id": 1}, {"id": 3}}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v but got %v", expected, result)
	} else {
		fmt.Println("TestF_Execute passed.")
	}
}

func TestF_Execute_All(t *testing.T) {
	// Test Case 1: No args, no kwargs
	f1 := NewF(func(v interface{}, args ...interface{}) interface{} {
		return v.(int) + 1
	})
	value1 := 10
	expected1 := 11
	result1, _ := f1.Execute(value1)
	if result1 != expected1 {
		t.Errorf("Expected %v but got %v", expected1, result1)
	}

	// Test Case 2: With kwargs
	f2 := NewF(func(v interface{}, args ...interface{}) interface{} {
		a := args[0].(string)
		b := args[1].(int)
		return fmt.Sprintf("%s%d", a, b+v.(int))
	}, "suffix", 100)
	value2 := 10
	expected2 := "suffix110"
	result2, _ := f2.Execute(value2)
	if result2 != expected2 {
		t.Errorf("Expected %v but got %v", expected2, result2)
	}

	// Test Case 3: With args
	f3 := NewF(func(v interface{}, args ...interface{}) interface{} {
		sum := v.(int)
		for _, arg := range args {
			sum += arg.(int)
		}
		return sum
	}, 20, 30, 40)
	value3 := 10
	expected3 := 100
	result3, _ := f3.Execute(value3)
	if result3 != expected3 {
		t.Errorf("Expected %v but got %v", expected3, result3)
	}
}
func TestF_Execute_No_args_No_kwargs(t *testing.T) {
	f1 := NewF(func(v interface{}, args ...interface{}) interface{} {
		return v.(int) + 1
	})
	value1 := 10
	expected1 := 11
	result1, _ := f1.Execute(value1)
	if result1 != expected1 {
		t.Errorf("Expected %v but got %v", expected1, result1)
	}
}

func TestF_Execute_with_kwargs(t *testing.T) {
	f2 := NewF(func(v interface{}, args ...interface{}) interface{} {
		suffix := args[0].(string)
		num := v.(int)
		return suffix + strconv.Itoa(num)
	}, "suffix")
	value2 := 10
	expected2 := "suffix10"
	result2, _ := f2.Execute(value2)
	if result2 != expected2 {
		t.Errorf("Expected %v but got %v", expected2, result2)
	}
}

func TestF_Execute_with_args(t *testing.T) {
	f3 := NewF(func(v interface{}, args ...interface{}) interface{} {
		sum := v.(int)
		for _, arg := range args {
			sum += arg.(int)
		}
		return sum
	}, 20, 30, 40)
	value3 := 10
	expected3 := 100
	result3, _ := f3.Execute(value3)
	if result3 != expected3 {
		t.Errorf("Expected %v but got %v", expected3, result3)
	}

}
