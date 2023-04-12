package whiteboard

import (
	"fmt"
	"reflect"
	"testing"
)

var testCases = []struct {
	name           string
	mapping        interface{}
	source         interface{}
	context        map[interface{}]interface{}
	expectedOutput interface{}
	expectedErr    error
}{
	{
		name:           "test case 1: empty mapping",
		mapping:        nil,
		source:         nil,
		context:        nil,
		expectedOutput: nil,
		expectedErr:    nil,
	},
	{
		name:           "test case 2: scalar value source",
		mapping:        "S(0)",
		source:         [2]string{"hello world", "world"},
		context:        nil,
		expectedOutput: "hello world",
		expectedErr:    nil,
	},
	{
		name:    "test case 3: list mapping",
		mapping: "S(\"a\",0,\"b\")",
		source: map[interface{}]interface{}{
			"a": []map[interface{}]interface{}{
				{"b": 42},
				{"c": 32},
			},
		},
		context:        nil,
		expectedOutput: 42,
		expectedErr:    nil,
	},
	{
		name:           "test case 4: scalar value with bender mapping",
		mapping:        "S(0) + \"!\"",
		source:         [2]string{"hello world", "world"},
		context:        nil,
		expectedOutput: "hello world!",
		expectedErr:    nil,
	},
	{
		name: "test case 5: complex mapping and non-scalar value source",
		mapping: map[interface{}]interface{}{
			"id":   "S(\"id\")",
			"name": "S(\"name\")",
			"rank": 42,
			"pets": "S(\"pets\")",
		},
		source: map[string]interface{}{
			"id":   123,
			"name": "Bob",
			"pets": []interface{}{
				map[string]interface{}{
					"name": "cat",
					"age":  2,
				},
				map[string]interface{}{
					"name": "dog",
					"age":  3,
				},
			},
		},
		context:        nil,
		expectedOutput: map[interface{}]interface{}{"id": 123, "name": "Bob", "pets": []interface{}{map[interface{}]interface{}{"age": 2, "name": "cat"}, map[interface{}]interface{}{"age": 3, "name": "dog"}}, "rank": 42},
		expectedErr:    nil,
	},
	{
		name: "test case 6: example with str and int",
		mapping: map[interface{}]interface{}{
			"level": "S(\"b\", \"userLevel\")",
			"kind":  "S(\"c\", \"userKind\")",
			"count": "K(1) + K(2)",
		},
		source: map[string]interface{}{
			"b": map[string]interface{}{
				"userLevel": 1,
				"name":      "123",
			},
			"c": map[string]interface{}{
				"userKind": "VIP",
				"age":      3,
			},
		},
		context:        nil,
		expectedOutput: map[interface{}]interface{}{"level": 1, "kind": "123", "count": 3},
		expectedErr:    nil,
	},
}

func TestBend_empty_mapping(t *testing.T) {
	tc := testCases[0]
	output, err := Bend(tc.mapping, tc.source, tc.context)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}
	if !reflect.DeepEqual(output, tc.expectedOutput) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}
}

func TestBend_empty_scalar_value(t *testing.T) {
	tc := testCases[1]
	fmt.Printf("%v\n", tc.name)
	output, err := Bend(tc.mapping, tc.source, tc.context)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}
	if !reflect.DeepEqual(output, tc.expectedOutput) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}
}

func TestBend_empty_list_mapping(t *testing.T) {
	tc := testCases[2]

	output, err := Bend(tc.mapping, tc.source, tc.context)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}
	if !reflect.DeepEqual(output, tc.expectedOutput) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}
}

func TestBend_empty_bender_complext(t *testing.T) {
	tc := testCases[4]
	output, err := Bend(tc.mapping, tc.source, tc.context)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}
	if mapsEqual(output.(map[interface{}]interface{}), tc.expectedOutput.(map[interface{}]interface{})) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}

}
func mapsEqual(map1, map2 map[interface{}]interface{}) bool {
	// 检查 map 的长度是否相等
	if len(map1) != len(map2) {
		return false
	}

	// 遍历第一个 map，检查其键值对是否在第二个 map 中都存在，并且对应的值相等
	for k, v1 := range map1 {
		v2, ok := map2[k]
		if !ok || !reflect.DeepEqual(v1, v2) {
			return false
		}
	}

	return true
}

func TestBend_empty_bender_test(t *testing.T) {
	tc := testCases[5]
	output, err := Bend(tc.mapping, tc.source, tc.context)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}
	if mapsEqual(output.(map[interface{}]interface{}), tc.expectedOutput.(map[interface{}]interface{})) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}

}
