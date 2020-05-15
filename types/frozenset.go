// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type FrozenSet struct {
	*Set
}

func NewFrozenSet() *FrozenSet {
	return &FrozenSet{NewSet()}
}

func NewFrozenSetFromSlice(slice []interface{}) *FrozenSet {
	return &FrozenSet{NewSetFromSlice(slice)}
}
