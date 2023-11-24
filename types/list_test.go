// Copyright 2023 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"bytes"
	"fmt"
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

func TestListStringer(t *testing.T) {
	pylist := func(sli ...interface{}) *List {
		return NewListFromSlice(sli)
	}

	lst := pylist(pylist(), pylist(1), pylist(2, 3), pylist(pylist(4, 5), 6))

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v", lst)
	got := buf.String()
	want := "[[], [1], [2, 3], [[4, 5], 6]]"

	if got != want {
		t.Fatalf("got= %q\nwant=%q", got, want)
	}
}
