// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

import (
	"encoding/binary"
	"math"
	"os"
)

type StorageClassInterface interface {
	New(size int, location string) StorageInterface
}

type StorageInterface interface {
	SetFromFile(f *os.File) error
}

type BaseStorage struct {
	Size     int
	Location string
}

// ----- Float -----

type FloatStorageClass struct{}
type FloatStorage struct {
	BaseStorage
	Data []float32
}

var _ StorageClassInterface = &FloatStorageClass{}
var _ StorageInterface = &FloatStorage{}

func (f *FloatStorageClass) New(size int, location string) StorageInterface {
	return &FloatStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

func (f *FloatStorage) SetFromFile(file *os.File) error {
	size, err := readSize(file)
	if err != nil {
		return err
	}

	data := make([]float32, size)
	buf := make([]byte, 512*4)
	for i, left := 0, 4*size; left > 0; {
		if left < len(buf) {
			buf = buf[0:left]
		}
		_, err := file.Read(buf)
		if err != nil {
			return err
		}
		for j := 0; j < len(buf); j += 4 {
			data[i] =
				math.Float32frombits(binary.LittleEndian.Uint32(buf[j : j+4]))
			i++
		}
		left -= len(buf)
	}
	f.Data = data
	return nil
}

func readSize(f *os.File) (int, error) {
	sizeBuf := make([]byte, 8)
	_, err := f.Read(sizeBuf)
	if err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint64(sizeBuf)), nil
}

// TODO: DoubleStorage
// TODO: HalfStorage
// TODO: LongStorage
// TODO: IntStorage
// TODO: ShortStorage
// TODO: CharStorage
// TODO: ByteStorage
// TODO: BoolStorage
// TODO: BFloat16Storage
// TODO: ComplexDoubleStorage
// TODO: ComplexFloatStorage
// TODO: QUInt8Storage
// TODO: QInt8Storage
// TODO: QInt32Storage
