package whiteboard

import (
	"errors"
	"reflect"
	"testing"
)

func TestIf_Execute(t *testing.T) {
	val1 := map[string]interface{}{
		"country":    "China",
		"first_name": "Li",
		"last_name":  "Na",
	}
	val2 := map[string]interface{}{
		"country":    "Brazil",
		"first_name": "Gustavo",
		"last_name":  "Kuerten",
	}

	testCases := []struct {
		name    string
		ifCond  Selector
		ifTrue  Selector
		ifFalse Selector
		input   interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:    "test condition true",
			ifCond:  NewExpressionSelector(&S{[]interface{}{"country"}}, &K{"China"}, "=="),
			ifTrue:  &S{[]interface{}{"first_name"}},
			ifFalse: &S{[]interface{}{"last_name"}},
			input:   val1,
			want:    "Li",
			wantErr: false,
		},
		{
			name:    "test condition false",
			ifCond:  NewExpressionSelector(&S{[]interface{}{"country"}}, &K{"China"}, "=="),
			ifTrue:  &S{[]interface{}{"first_name"}},
			ifFalse: &S{[]interface{}{"last_name"}},
			input:   val2,
			want:    "Kuerten",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ifStmt := NewIf(tc.ifCond, tc.ifTrue, tc.ifFalse)
			got, err := ifStmt.Execute(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ifStmt.Execute() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ifStmt.Execute() got = %v, want %v", got, tc.want)
			}
		})
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
