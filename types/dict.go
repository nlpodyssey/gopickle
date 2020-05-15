// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type DictSetter interface {
	Set(k, v interface{})
}

type Dict map[interface{}]interface{}

var _ DictSetter = &Dict{}

func NewDict() *Dict {
	d := make(Dict)
	return &d
}

func (d *Dict) Set(k, v interface{}) {
	(*d)[k] = v
}

func (d *Dict) Get(k interface{}) (interface{}, bool) {
	value, ok := (*d)[k]
	return value, ok
}

func (d *Dict) Len() int {
	return len(*d)
}
