package whiteboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/ghodss/yaml"
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
		expectedErr:    errors.New("mapping or source is empty"),
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
		expectedOutput: map[interface{}]interface{}{"level": 1, "kind": "VIP", "count": 3},
		expectedErr:    nil,
	},

	{
		name: "test case 7: example with control flow: IF",
		mapping: map[string]interface{}{
			"name": "IF( ExpS( S(\"country\"), K(\"China\"), \"==\" ), S(\"first_name\"), S(\"last_name\"))",
		},
		source: map[string]interface{}{
			"country":    "China",
			"first_name": "Li",
			"last_name":  "Na",
		},
		context:        nil,
		expectedOutput: map[string]interface{}{"name": "Li"},
		expectedErr:    nil,
	},

	{
		name: "test case 8: example with control flow:Alternation",
		mapping: map[string]string{
			"name": "AL( S(1), S(0), S(\"key1\"))",
		},
		source:         []interface{}{"a", "b"},
		context:        nil,
		expectedOutput: map[string]interface{}{"name": "b"},
		expectedErr:    nil,
	},
}

var ActionMaps = map[interface{}]map[interface{}]interface{}{
	"init": {
		"userId": "123",
	},
	"a": {
		"userAgent": "agent",
		"userName":  "name",
	},
	"b": {
		"userLevel": "level",
	},
	"c": {
		"userKind": "kind",
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

func TestBend_empty_bender_mapping(t *testing.T) {
	tc := testCases[3]

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

	fmt.Println(output)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}
	expect := tc.expectedOutput

	if !CompareMaps(output, expect) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}
}

func CompareMaps(mapA interface{}, mapB interface{}) bool {
	// 判断类型是否为map
	if reflect.TypeOf(mapA).Kind() != reflect.Map || reflect.TypeOf(mapB).Kind() != reflect.Map {
		return false
	}

	valueA := reflect.ValueOf(mapA)
	valueB := reflect.ValueOf(mapB)

	// 获取两个map的键名集合并排序
	keysA := valueA.MapKeys()
	keysB := valueB.MapKeys()
	sort.Slice(keysA, func(i, j int) bool {
		return keysA[i].String() < keysA[j].String()
	})
	sort.Slice(keysB, func(i, j int) bool {
		return keysB[i].String() < keysB[j].String()
	})

	// 如果两个map的键名数量不同，则直接返回false
	if len(keysA) != len(keysB) {
		return false
	}

	// 遍历第一个map，检查其键和值是否都存在于第二个map中
	for i, key := range keysA {
		valueA := valueA.MapIndex(key)
		valueB := valueB.MapIndex(keysB[i])

		if !valueB.IsValid() || !reflect.DeepEqual(valueA.Interface(), valueB.Interface()) {
			fmt.Printf("%v-->%v / %v\n", key, valueA.Interface(), valueB.Interface())
			return false
		}
	}

	// 如果两个map内的数据完全一致，则返回true
	return true
}

func TestBend_empty_bender_test(t *testing.T) {
	tc := testCases[5]
	fmt.Println(tc.name)
	output, err := Bend(tc.mapping, tc.source, tc.context)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}

	expect := tc.expectedOutput

	if !CompareMaps(output, expect) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}

}

func TestBend_testing(t *testing.T) {
	//读取YAML内容
	jsonBytes := readYAML()

	var data map[string]interface{}
	json.Unmarshal([]byte(jsonBytes), &data)

	fmt.Println("-----")
	stars, ok := data["spec"].(map[string]interface{})["stars"].([]interface{})
	if !ok {
		fmt.Println("Failed to get stars")
		return
	}

	for _, star := range stars {
		starMap, ok := star.(map[string]interface{})
		if !ok {
			fmt.Println("Failed to get star map")
			continue
		}
		fmt.Println(starMap["name"], starMap["image"], starMap["port"], starMap["action"], starMap["dependencies"], starMap["param"])

		parms, err := yaml.YAMLToJSON([]byte(starMap["param"].(string)))

		if err != nil {
			panic(err)
		}

		var mapping map[string]interface{}
		json.Unmarshal([]byte(parms), &mapping)

		fmt.Println(mapping)
		output, err := Bend(mapping, ActionMaps, nil)

		fmt.Printf("%v-->%v\n", starMap["name"], output)
		if err != nil {
			panic(err)
		}

		fmt.Println("-----")
	}

}

func readYAML() []byte {
	yamlFile, err := os.Open("samples/astro_new.yaml")
	if err != nil {
		panic(err)
	}
	defer yamlFile.Close()

	yamlBytes, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		panic(err)
	}

	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))
	return jsonBytes
}

func TestBend_with_IF(t *testing.T) {
	tc := testCases[6]
	fmt.Printf("%v-->%v-->%v\n", tc.name, tc.mapping, tc.source)

	output, err := Bend(tc.mapping, tc.source, tc.context)

	fmt.Println(output)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}

	if !CompareMaps(output, tc.expectedOutput) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}

}

func TestBend_with_Al(t *testing.T) {
	tc := testCases[7]
	fmt.Printf("%v-->%v-->%v\n", tc.name, tc.mapping, tc.source)

	output, err := Bend(tc.mapping, tc.source, tc.context)

	fmt.Println(output)
	if !reflect.DeepEqual(err, tc.expectedErr) {
		t.Errorf("expected error %v, but got %v", tc.expectedErr, err)
	}

	if !CompareMaps(output, tc.expectedOutput) {
		t.Errorf("expected output %v, but got %v", tc.expectedOutput, output)
	}
}

func Test_2(t *testing.T) {
	param := make(map[string]string)
	var ParamFormat = "{\"id\":\"S(\\\"init\\\",\\\"userId\\\")\"}"
	var ActionMaps = map[string]map[string]interface{}{
		"init": {
			"userId": "123",
		},
	}
	json.Unmarshal([]byte(ParamFormat), &param)
	fmt.Println("ActionMap: ", ActionMaps)
	fmt.Println("node param format: ", ParamFormat)
	fmt.Println("rule param:", param)
	data, err := Bend(param, ActionMaps, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("node result data:", data)
	expect := map[string]interface{}{
		"id": "123",
	}
	if !CompareMaps(data, expect) {
		t.Errorf("expected output %v, but got %v", expect, data)
	}
}
