// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"strings"
)

// ListAppender is implemented by any value that exhibits a list-like
// behaviour, allowing arbitrary values to be appended.
type ListAppender interface {
	Append(v interface{})
}

// List represents a Python "list" (builtin type).
type List []interface{}

var _ ListAppender = &List{}

// NewList makes and returns a new empty List.
func NewList() *List {
	l := make(List, 0, 4)
	return &l
}

// NewListFromSlice makes and returns a new List initialized with the elements
// of the given slice.
//
// The new List is a simple type cast of the input slice; the slice is _not_
// copied.
func NewListFromSlice(slice []interface{}) *List {
	l := List(slice)
	return &l
}

// Append appends one element to the end of the List.
func (l *List) Append(v interface{}) {
	*l = append(*l, v)
}

// Get returns the element of the List at the given index.
//
// It panics if the index is out of range.
func (l *List) Get(i int) interface{} {
	return (*l)[i]
}

// Len returns the length of the List.
func (l *List) Len() int {
	return len(*l)
}

func (*List) Call(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return NewList(), nil
	}
	if len(args) == 1 {
		return args[0], nil
	}
	return nil, fmt.Errorf("List: invalid arguments: %#v", args)
}

func (l *List) String() string {
	if l == nil {
		return "nil"
	}
	o := new(strings.Builder)
	o.WriteString("[")
	for i, v := range *l {
		if i > 0 {
			o.WriteString(", ")
		}
		fmt.Fprintf(o, "%v", v)
	}
	o.WriteString("]")
	return o.String()
}
