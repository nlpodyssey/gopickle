// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

type Tensor struct {
	Source        StorageInterface
	StorageOffset int
	Size          []int
	Stride        []int
	RequiresGrad  bool
}
