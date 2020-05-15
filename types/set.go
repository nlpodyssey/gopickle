// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type SetAdder interface {
	Add(v interface{})
}

type Set map[interface{}]setEmptyStruct

var _ SetAdder = &Set{}

type setEmptyStruct struct{}

func NewSet() *Set {
	s := make(Set)
	return &s
}

func NewSetFromSlice(slice []interface{}) *Set {
	s := make(Set, len(slice))
	for _, item := range slice {
		s[item] = setEmptyStruct{}
	}
	return &s
}

func (s *Set) Len() int {
	return len(*s)
}

func (s *Set) Add(v interface{}) {
	(*s)[v] = setEmptyStruct{}
}

func (s *Set) Has(v interface{}) bool {
	_, ok := (*s)[v]
	return ok
}
