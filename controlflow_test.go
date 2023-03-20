package whiteboard

// func TestIf_Execute(t *testing.T) {
// 	// test case 1
// 	condition1 := func(val interface{}) bool {
// 		return val.(map[string]interface{})["country"].(string) == "China"
// 	}
// 	whenTrue1 := NewS("first_name")
// 	whenFalse1 := NewS("last_name")
// 	if_1 := NewIf(condition1, whenTrue1, whenFalse1)
// 	val1 := map[string]interface{}{"country": "China", "first_name": "Li", "last_name": "Na"}
// 	res1 := if_1.Execute(val1)
// 	if res1 != "Li" {
// 		t.Errorf("Test case 1 failed: expected %v but got %v", "Li", res1)
// 	}

// 	// test case 2
// 	condition2 := func(val interface{}) bool {
// 		return val.(map[string]interface{})["age"].(int) < 18
// 	}
// 	whenTrue2 := NewS("You are not allowed to vote yet.")
// 	whenFalse2 := NewS("You can vote now.")
// 	if_2 := NewIf(condition2, whenTrue2, whenFalse2)
// 	val2 := map[string]interface{}{"name": "Alice", "age": 21}
// 	res2 := if_2.Execute(val2)
// 	if res2 != "You can vote now." {
// 		t.Errorf("Test case 2 failed: expected %v but got %v", "You can vote now.", res2)
// 	}
// }
