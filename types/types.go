// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type Callable interface {
	Call(args ...interface{}) (interface{}, error)
}

type PyNewable interface {
	// __new__
	PyNew(args ...interface{}) (interface{}, error)
}

type PyStateSettable interface {
	// __setstate__
	PySetState(state interface{}) error
}

type PyDictSettable interface {
	// __dict__
	PyDictSet(key, value interface{}) error
}

type PyAttrSettable interface {
	// setattr()
	PySetAttr(key string, value interface{}) error
}
