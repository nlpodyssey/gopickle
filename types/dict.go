// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"reflect"
)

// DictSetter is implemented by any value that exhibits a dict-like behaviour,
// allowing arbitrary key/value pairs to be set.
type DictSetter interface {
	Set(key, value interface{})
}

// Dict represents a Python "dict" (builtin type).
//
// It is implemented as a slice, instead of a map, because in Go not
// all types can be map's keys (e.g. slices).
type Dict []*DictEntry

type DictEntry struct {
	Key   interface{}
	Value interface{}
}

var _ DictSetter = &Dict{}

// NewDict makes and returns a new empty Dict.
func NewDict() *Dict {
	d := make(Dict, 0)
	return &d
}

// Set sets into the Dict the given key/value pair.
func (d *Dict) Set(key, value interface{}) {
	*d = append(*d, &DictEntry{
		Key:   key,
		Value: value,
	})
}

// Get returns the value associated with the given key (if any), and whether
// the key is present or not.
func (d *Dict) Get(key interface{}) (interface{}, bool) {
	for _, entry := range *d {
		if reflect.DeepEqual(entry.Key, key) {
			return entry.Value, true
		}
	}
	return nil, false
}

// Len returns the length of the Dict, that is, the amount of key/value pairs
// contained by the Dict.
func (d *Dict) Len() int {
	return len(*d)
}
