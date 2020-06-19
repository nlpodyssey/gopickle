// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// DictSetter is implemented by any value that exhibits a dict-like behaviour,
// allowing arbitrary key/value pairs to be set.
type DictSetter interface {
	Set(key, value interface{})
}

// Dict represents a Python "dict" (builtin type).
type Dict map[interface{}]interface{}

var _ DictSetter = &Dict{}

// NewDict makes and returns a new empty Dict.
func NewDict() *Dict {
	d := make(Dict)
	return &d
}

// Set sets into the Dict the given key/value pair.
//
// If the key is already present, its associated value is replaced with the
// new one.
func (d *Dict) Set(key, value interface{}) {
	(*d)[key] = value
}

// Get returns the value associated with the given key (if any), and whether
// the key is present or not.
func (d *Dict) Get(key interface{}) (interface{}, bool) {
	value, ok := (*d)[key]
	return value, ok
}

// Len returns the length of the Dict, that is, the amount of key/value pairs
// contained by the Dict.
func (d *Dict) Len() int {
	return len(*d)
}
