// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"container/list"
	"fmt"
)

type OrderedDictClass struct{}

var _ Callable = &OrderedDictClass{}

func (*OrderedDictClass) Call(args ...interface{}) (interface{}, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf(
			"OrderedDictClass.Call args not supported: %#v", args)
	}
	return NewOrderedDict(), nil
}

type OrderedDict struct {
	Map    map[interface{}]*OrderedDictEntry
	List   *list.List
	PyDict map[string]interface{} // __dict__
}

var _ DictSetter = &OrderedDict{}
var _ PyDictSettable = &OrderedDict{}

type OrderedDictEntry struct {
	Key         interface{}
	Value       interface{}
	ListElement *list.Element
}

func NewOrderedDict() *OrderedDict {
	return &OrderedDict{
		Map:    make(map[interface{}]*OrderedDictEntry),
		List:   list.New(),
		PyDict: make(map[string]interface{}),
	}
}

func (o *OrderedDict) Set(k, v interface{}) {
	if entry, ok := o.Map[k]; ok {
		entry.Value = v
		return
	}

	entry := &OrderedDictEntry{
		Key:   k,
		Value: v,
	}
	entry.ListElement = o.List.PushBack(entry)
	o.Map[k] = entry
}

func (o *OrderedDict) Get(k interface{}) (interface{}, bool) {
	entry, ok := o.Map[k]
	if !ok {
		return nil, false
	}
	return entry.Value, true
}

func (o *OrderedDict) Len() int {
	return len(o.Map)
}

func (o *OrderedDict) PyDictSet(key, value interface{}) error {
	sKey, keyOk := key.(string)
	if !keyOk {
		return fmt.Errorf(
			"OrderedDict.PyDictSet() requires string key: %#v", key)
	}
	o.PyDict[sKey] = value
	return nil
}
