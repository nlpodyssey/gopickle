// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"reflect"
	"strings"
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
type Dict []DictEntry

type DictEntry struct {
	Key   interface{}
	Value interface{}
}

var _ DictSetter = &Dict{}

// NewDict makes and returns a new empty Dict.
func NewDict() *Dict {
	d := make(Dict, 0, 4)
	return &d
}

// Set sets into the Dict the given key/value pair.
func (d *Dict) Set(key, value interface{}) {
	*d = append(*d, DictEntry{
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

// MustGet returns the value associated with the given key, if if it exists,
// otherwise it panics.
func (d *Dict) MustGet(key interface{}) interface{} {
	value, ok := d.Get(key)
	if !ok {
		panic(fmt.Errorf("key not found in Dict: %#v", key))
	}
	return value
}

// Len returns the length of the Dict, that is, the amount of key/value pairs
// contained by the Dict.
func (d *Dict) Len() int {
	return len(*d)
}

// Keys returns the keys of the dict
func (d *Dict) Keys() []interface{} {
	out := make([]interface{}, len(*d))
	for i, entry := range *d {
		out[i] = entry.Key
	}

	return out
}

func (*Dict) Call(args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return NewDict(), nil
	}
	if len(args) == 1 {
		return args[0], nil
	}
	return nil, fmt.Errorf("Dict: invalid arguments: %#v", args)
}

func (d *Dict) String() string {
	if d == nil {
		return "nil"
	}
	o := new(strings.Builder)
	o.WriteString("{")
	for i, e := range *d {
		if i > 0 {
			o.WriteString(", ")
		}
		fmt.Fprintf(o, "%v: %v", e.Key, e.Value)
	}
	o.WriteString("}")
	return o.String()
}
