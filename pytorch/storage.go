// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

import (
	"encoding/binary"
	"io"
	"math"
)

type StorageClassInterface interface {
	New(size int, location string) StorageInterface
}

type StorageInterface interface {
	SetFromFile(r io.Reader) error
	SetFromFileWithSize(r io.Reader, size int) error
}

type BaseStorage struct {
	Size     int
	Location string
}

// ----- Half -----

type HalfStorageClass struct{}

var _ StorageClassInterface = &HalfStorageClass{}

func (f *HalfStorageClass) New(size int, location string) StorageInterface {
	return &HalfStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type HalfStorage struct {
	BaseStorage
	Data []float32
}

var _ StorageInterface = &HalfStorage{}

func (f *HalfStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *HalfStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]float32, size)
	br := NewLimitedBufferReader(r, size, 2, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		u16 := binary.LittleEndian.Uint16(bytes)
		data[i] = math.Float32frombits(FloatBits16to32(u16))
	}
	f.Data = data
	return nil
}

// ----- Float -----

type FloatStorageClass struct{}

var _ StorageClassInterface = &FloatStorageClass{}

func (f *FloatStorageClass) New(size int, location string) StorageInterface {
	return &FloatStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type FloatStorage struct {
	BaseStorage
	Data []float32
}

var _ StorageInterface = &FloatStorage{}

func (f *FloatStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *FloatStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]float32, size)
	br := NewLimitedBufferReader(r, size, 4, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		data[i] = math.Float32frombits(binary.LittleEndian.Uint32(bytes))
	}
	f.Data = data
	return nil
}

// ----- Double -----

type DoubleStorageClass struct{}

var _ StorageClassInterface = &DoubleStorageClass{}

func (f *DoubleStorageClass) New(size int, location string) StorageInterface {
	return &DoubleStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type DoubleStorage struct {
	BaseStorage
	Data []float64
}

var _ StorageInterface = &DoubleStorage{}

func (f *DoubleStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *DoubleStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]float64, size)
	br := NewLimitedBufferReader(r, size, 8, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		data[i] = math.Float64frombits(binary.LittleEndian.Uint64(bytes))
	}
	f.Data = data
	return nil
}

// ----- Char -----

type CharStorageClass struct{}

var _ StorageClassInterface = &CharStorageClass{}

func (f *CharStorageClass) New(size int, location string) StorageInterface {
	return &CharStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type CharStorage struct {
	BaseStorage
	Data []int8
}

var _ StorageInterface = &CharStorage{}

func (f *CharStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *CharStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]int8, size)
	br := NewLimitedBufferReader(r, size, 1, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		data[i] = int8(bytes[0])
	}
	f.Data = data
	return nil
}

// ----- Short -----

type ShortStorageClass struct{}

var _ StorageClassInterface = &ShortStorageClass{}

func (f *ShortStorageClass) New(size int, location string) StorageInterface {
	return &ShortStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type ShortStorage struct {
	BaseStorage
	Data []int16
}

var _ StorageInterface = &ShortStorage{}

func (f *ShortStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *ShortStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]int16, size)
	br := NewLimitedBufferReader(r, size, 2, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		data[i] = int16(binary.LittleEndian.Uint16(bytes))
	}
	f.Data = data
	return nil
}

// ----- Int -----

type IntStorageClass struct{}

var _ StorageClassInterface = &IntStorageClass{}

func (f *IntStorageClass) New(size int, location string) StorageInterface {
	return &IntStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type IntStorage struct {
	BaseStorage
	Data []int32
}

var _ StorageInterface = &IntStorage{}

func (f *IntStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *IntStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]int32, size)
	br := NewLimitedBufferReader(r, size, 4, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		data[i] = int32(binary.LittleEndian.Uint32(bytes))
	}
	f.Data = data
	return nil
}

// ----- Long -----

type LongStorageClass struct{}

var _ StorageClassInterface = &LongStorageClass{}

func (f *LongStorageClass) New(size int, location string) StorageInterface {
	return &LongStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type LongStorage struct {
	BaseStorage
	Data []int64
}

var _ StorageInterface = &LongStorage{}

func (f *LongStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *LongStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]int64, size)
	br := NewLimitedBufferReader(r, size, 8, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		data[i] = int64(binary.LittleEndian.Uint64(bytes))
	}
	f.Data = data
	return nil
}

// ----- Byte -----

type ByteStorageClass struct{}

var _ StorageClassInterface = &ByteStorageClass{}

func (f *ByteStorageClass) New(size int, location string) StorageInterface {
	return &ByteStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type ByteStorage struct {
	BaseStorage
	Data []uint8
}

var _ StorageInterface = &ByteStorage{}

func (f *ByteStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *ByteStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]uint8, size)
	br := NewLimitedBufferReader(r, size, 1, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		data[i] = bytes[0]
	}
	f.Data = data
	return nil
}

// ----- Bool -----

type BoolStorageClass struct{}

var _ StorageClassInterface = &BoolStorageClass{}

func (f *BoolStorageClass) New(size int, location string) StorageInterface {
	return &BoolStorage{
		BaseStorage: BaseStorage{Size: size, Location: location},
		Data:        nil,
	}
}

type BoolStorage struct {
	BaseStorage
	Data []bool
}

var _ StorageInterface = &BoolStorage{}

func (f *BoolStorage) SetFromFile(r io.Reader) error {
	return setFromFile(f, r)
}

func (f *BoolStorage) SetFromFileWithSize(r io.Reader, size int) error {
	data := make([]bool, size)
	br := NewLimitedBufferReader(r, size, 1, 512)
	for i := 0; i < size; i++ {
		bytes, err := br.ReadNext()
		if err != nil {
			return err
		}
		data[i] = bytes[0] == 1
	}
	f.Data = data
	return nil
}

func setFromFile(s StorageInterface, r io.Reader) error {
	sizeBuf := make([]byte, 8)
	_, err := r.Read(sizeBuf)
	if err != nil {
		return err
	}
	size := int(binary.LittleEndian.Uint64(sizeBuf))
	return s.SetFromFileWithSize(r, size)
}
