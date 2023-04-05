package whiteboard

import (
	"errors"
	"testing"
)

func TestIf_Execute_str(t *testing.T) {
	// test case 1
	condition1 := func(val interface{}) bool {
		return val.(map[string]interface{})["country"].(string) == "China"
	}
	whenTrue1, _ := NewS("first_name")
	whenFalse1, _ := NewS("last_name")
	if_1 := NewIf(condition1, whenTrue1, whenFalse1)
	val1 := map[string]interface{}{"country": "China", "first_name": "Li", "last_name": "Na"}
	res1, _ := if_1.Execute(val1)
	if res1 != "Li" {
		t.Errorf("Test case 1 failed: expected %v but got %v", "Li", res1)
	}

}

func TestIf_Execute_int(t *testing.T) {
	// test case 2
	condition2 := func(val interface{}) bool {
		return val.(map[string]interface{})["age"].(int) < 18
	}
	whenTrue2, _ := NewK("You are not allowed to vote yet.")
	whenFalse2, _ := NewK("You can vote now.")
	if_2 := NewIf(condition2, whenTrue2, whenFalse2)
	val2 := map[string]interface{}{"name": "Alice", "age": 21}
	res2, err := if_2.Execute(val2)
	if err != nil {
		t.Errorf("Unexpected error returned: %v", err)
	}
	if res2 != "You can vote now." {
		t.Errorf("Test case 2 failed: expected %v but got %v", "You can vote now.", res2)
	}
}

func TestAlternation_Execute_list(t *testing.T) {
	// b := NewAlternation(
	// 	NewS(1), NewS(0), NewS("key1"),
	// )

	s1, _ := NewS(1)
	s2, _ := NewS(0)
	s3, _ := NewS("key1")
	b := NewAlternation(
		s1, s2, s3,
	)

	// Test cases for list inputs
	lst1 := []interface{}{"a", "b"}
	expected1 := "b"
	if res, err := b.Execute(lst1); res != expected1 {
		t.Errorf("Unexpected result: %v (expected %v) error:%v", res, expected1, err)
	}

	lst2 := []interface{}{"a"}
	expected2 := "a"
	if res, err := b.Execute(lst2); res != expected2 {
		t.Errorf("Unexpected result: %v (expected %v) error:%v", res, expected2, err)
	}

	lst3 := []interface{}{}
	expected3 := errors.New("KeyError")
	if _, err := b.Execute(lst3); err == nil {
		t.Errorf("Unexpected error: %v (expected %v)", err, expected3)
	}

}

func TestAlternation_Execute_dict(t *testing.T) {
	s1, _ := NewS(1)
	s2, _ := NewS(0)
	s3, _ := NewS("key1")
	s4, _ := NewS("other_key")
	b := NewAlternation(
		s1, s2, s3, s4,
	)

	// Test cases for dictionary inputs
	dict1 := map[string]interface{}{
		"key1": 23,
	}
	expected4 := 23
	if res, _ := b.Execute(dict1); res != expected4 {
		t.Errorf("Unexpected result: %v (expected %v)", res, expected4)
	}

	dict2 := map[string]interface{}{
		"other_key": "value",
	}
	expected5 := "value"
	if res, _ := b.Execute(dict2); res != expected5 {
		t.Errorf("Unexpected result: %v (expected %v)", res, expected5)
	}

	dict3 := map[string]interface{}{}
	expected6 := errors.New("KeyError")
	if _, err := b.Execute(dict3); err == nil {
		t.Errorf("Unexpected error: %v (expected %v)", err, expected6)
	}
}

/**
func TestSwitch_Execute_Example(t *testing.T) {
	// create Switch object with cases for 'twitter' and 'mastodon', and default case for 'email'
	ss, _ := NewS("service")
	sh, _ = NewS("handle")
	k, _ := NewK("@")
	s, _ := NewS("server")

	se, _ := NewS("email")

	switchSelector := &Switch{
		keySelctor: ss,
		cases: map[interface{}]selector{
			"twitter":  sh,
			"mastodon": sh + k + s,
		},
		defaultSelector: se,
	}

	// test case for 'twitter' service
	source := map[interface{}]interface{}{"service": "twitter", "handle": "etandel"}
	expectedResult := "etandel"
	result, err := switchSelector.Execute(source)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected result to be %v but got %v", expectedResult, result)
	}

	// test case for 'mastodon' service
	source = map[interface{}]interface{}{"service": "mastodon", "handle": "etandel", "server": "mastodon.social"}
	expectedResult = "etandel@mastodon.social"
	result, err = switchSelector.Execute(source)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected result to be %v but got %v", expectedResult, result)
	}

	// test default case for 'facebook' service
	source = map[interface{}]interface{}{"service": "facebook", "email": "email@whatever.com"}
	expectedResult = "email@whatever.com"
	result, err = switchSelector.Execute(source)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != expectedResult {
		t.Errorf("Expected result to be %v but got %v", expectedResult, result)
	}
}


**/
