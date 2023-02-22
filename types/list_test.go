package types

import (
	"testing"
)

func TestCall(t *testing.T) {
	list := NewList()
	list.Append("foo")
	args := []interface{}{list}
	result, _ := list.Call(args)
	actual := (*result.([]interface{})[0].(*List))[0]
	expected := "foo"
	if actual != expected {
		t.Errorf("expected %v, actual: %v", expected, actual)
	}
}
