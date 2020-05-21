// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type ByteArray []byte

func NewByteArray() *ByteArray {
	b := make(ByteArray, 0)
	return &b
}

func NewByteArrayFromSlice(slice []byte) *ByteArray {
	b := ByteArray(slice)
	return &b
}

func (b *ByteArray) Get(i int) byte {
	return (*b)[i]
}

func (b *ByteArray) Len() int {
	return len(*b)
}
