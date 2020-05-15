// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

import (
	"fmt"
	"github.com/nlpodyssey/gopickle/types"
)

type RebuildTensorV2 struct{}

var _ types.Callable = &RebuildTensorV2{}

func (r *RebuildTensorV2) Call(args ...interface{}) (interface{}, error) {
	if len(args) != 6 {
		return nil, fmt.Errorf("RebuildTensorV2 unexpected args: %#v", args)
	}
	storage, storageOk := args[0].(StorageInterface)
	storageOffset, storageOffsetOk := args[1].(int)
	size, sizeOk := args[2].(*types.Tuple)
	stride, strideOk := args[3].(*types.Tuple)
	requiresGrad, requiresGradOk := args[4].(bool)
	// arg[5] "backward hooks" is unused
	if !storageOk || !storageOffsetOk || !sizeOk || !strideOk ||
		!requiresGradOk {
		return nil, fmt.Errorf("RebuildTensorV2 unexpected args: %#v", args)
	}

	tensor := &Tensor{
		Source:        storage,
		StorageOffset: storageOffset,
		RequiresGrad:  requiresGrad,
	}
	var err error
	tensor.Size, err = tupleToIntSlice(size)
	if err != nil {
		return nil, err
	}
	tensor.Stride, err = tupleToIntSlice(stride)
	if err != nil {
		return nil, err
	}
	return tensor, nil
}

func tupleToIntSlice(tuple *types.Tuple) ([]int, error) {
	length := tuple.Len()
	slice := make([]int, length)
	for i := 0; i < length; i++ {
		value, ok := tuple.Get(i).(int)
		if !ok {
			return nil, fmt.Errorf("tuple of ints expected: %#v", tuple)
		}
		slice[i] = value
	}
	return slice, nil
}
