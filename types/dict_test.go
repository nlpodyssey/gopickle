package types

import "testing"

func TestDictCall(t *testing.T) {
	d := NewDict()
	d.Set("foo", "bar")
	args := []interface{}{d}
	result, _ := d.Call(args)
	resultdict := *result.([]interface{})[0].(*Dict)
	actual, _ := resultdict.Get("foo")
	expected := "bar"
	if actual != expected {
		t.Errorf("expected %v, actual: %v", expected, actual)
	}
}
