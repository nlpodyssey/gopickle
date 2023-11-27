// Copyright 2023 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"bytes"
	"fmt"
	"testing"
)

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

func TestDictStringer(t *testing.T) {
	dct := NewDict()
	dct.Set("empty", NewDict())
	dct.Set("one", "un")
	dct.Set("two", "deux")
	sub := NewDict()
	sub.Set("eins", "one")
	sub.Set("zwei", []string{"two", "deux"})
	sub.Set(2, []int{2 * 2, 2 * 3, 2 * 4})
	dct.Set("sub", sub)

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v", dct)

	var (
		got  = buf.String()
		want = "{empty: {}, one: un, two: deux, sub: {eins: one, zwei: [two deux], 2: [4 6 8]}}"
	)

	if got != want {
		t.Fatalf("got= %q\nwant=%q", got, want)
	}
}
