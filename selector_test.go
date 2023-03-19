package whiteboard

import (
	"reflect"
	"sort"
	"testing"
)

func TestK_Execute(t *testing.T) {
	k := NewK("Hello World")
	result := k.Execute(nil)
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', but got %v", result)
	}
}

func TestS_Execute(t *testing.T) {
	data := map[interface{}]interface{}{
		"a": []map[interface{}]interface{}{
			{"b": 42},
		},
	}

	s, err := NewS("a", 0, "b")
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}

	result := s.Execute(data)
	if result != 42 {
		t.Errorf("Expected 42, but got %v", result)
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
	var data = []map[string]interface{}{
		{"id": 3},
		{"id": 1},
	}

	// f := NewF(
	// 	func(value interface{}, key string) interface{} {
	// 		list := value.([]map[string]interface{})
	// 		sort.Slice(list, func(i, j int) bool {
	// 			return list[i][key].(int) < list[j][key].(int)
	// 		})
	// 		return list
	// 	},
	// 	"key", "id",
	// )

	f := NewF(func(value interface{}, args ...interface{}) interface{} {
		return mySortFunc(append([]interface{}{value}, args...)...)
	}, "key", "id")

	result := f.Execute(data)
	expected := []map[string]interface{}{
		{"id": 1},
		{"id": 3},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
