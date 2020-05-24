// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type ListAppender interface {
	Append(v interface{})
}

type List []interface{}

var _ ListAppender = &List{}

func NewList() *List {
	l := make(List, 0)
	return &l
}

func NewListFromSlice(slice []interface{}) *List {
	l := List(slice)
	return &l
}

func (l *List) Append(v interface{}) {
	*l = append(*l, v)
}

func (l *List) Get(i int) interface{} {
	return (*l)[i]
}

func (l *List) Len() int {
	return len(*l)
}
