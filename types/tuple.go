// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"strings"
)

type Tuple []interface{}

func NewTupleFromSlice(slice []interface{}) *Tuple {
	t := Tuple(slice)
	return &t
}

func (t *Tuple) Get(i int) interface{} {
	return (*t)[i]
}

func (t *Tuple) Len() int {
	return len(*t)
}

func (t *Tuple) String() string {
	if t == nil {
		return "nil"
	}
	o := new(strings.Builder)
	o.WriteString("(")
	for i, v := range *t {
		if i > 0 {
			o.WriteString(", ")
		}
		fmt.Fprintf(o, "%v", v)
	}
	if t.Len() == 1 {
		o.WriteString(",")
	}
	o.WriteString(")")
	return o.String()
}
