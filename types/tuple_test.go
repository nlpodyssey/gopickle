// Copyright 2023 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTupleStringer(t *testing.T) {
	tuple := func(sli ...interface{}) *Tuple {
		return NewTupleFromSlice(sli)
	}

	tup := tuple(tuple(), tuple(1), tuple(2, 3), tuple(tuple(4, 5), 6))

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v", tup)
	got := buf.String()
	want := "((), (1,), (2, 3), ((4, 5), 6))"

	if got != want {
		t.Fatalf("got= %q\nwant=%q", got, want)
	}
}
