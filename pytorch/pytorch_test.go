// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

import (
	"fmt"
	"path"
	"testing"
)

func TestFloat16Tensors(t *testing.T) { // Half
	for _, filename := range makeFilenames("tensor_float16") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*HalfStorage)
			if !fsOk {
				t.Fatalf("expected *HalfStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertFloat32SliceEqual(t, fs.Data, []float32{1.2, -3.4, 5.6, -7.8}, 0.002)
		})
	}
}

func TestFloat32Tensors(t *testing.T) { // Float
	for _, filename := range makeFilenames("tensor_float32") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*FloatStorage)
			if !fsOk {
				t.Fatalf("expected *FloatStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertFloat32SliceEqual(t, fs.Data, []float32{1.2, -3.4, 5.6, -7.8}, 0.0)
		})
	}
}

func TestFloat64Tensors(t *testing.T) { // Double
	for _, filename := range makeFilenames("tensor_float64") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*DoubleStorage)
			if !fsOk {
				t.Fatalf("expected *DoubleStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertFloat64SliceEqual(t, fs.Data, []float64{1.2, -3.4, 5.6, -7.8}, 0.0)
		})
	}
}

func TestInt8Tensors(t *testing.T) { // Char
	for _, filename := range makeFilenames("tensor_int8") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*CharStorage)
			if !fsOk {
				t.Fatalf("expected *CharStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertInt8SliceEqual(t, fs.Data, []int8{1, -2, 3, -4})
		})
	}
}

func TestInt16Tensors(t *testing.T) { // Short
	for _, filename := range makeFilenames("tensor_int16") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*ShortStorage)
			if !fsOk {
				t.Fatalf("expected *ShortStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertInt16SliceEqual(t, fs.Data, []int16{1, -2, 3, -4})
		})
	}
}

func TestInt32Tensors(t *testing.T) { // Int
	for _, filename := range makeFilenames("tensor_int32") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*IntStorage)
			if !fsOk {
				t.Fatalf("expected *IntStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertInt32SliceEqual(t, fs.Data, []int32{1, -2, 3, -4})
		})
	}
}

func TestInt64Tensors(t *testing.T) { // Long
	for _, filename := range makeFilenames("tensor_int64") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*LongStorage)
			if !fsOk {
				t.Fatalf("expected *LongStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertInt64SliceEqual(t, fs.Data, []int64{1, -2, 3, -4})
		})
	}
}

func TestUInt8Tensors(t *testing.T) { // Byte
	for _, filename := range makeFilenames("tensor_uint8") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*ByteStorage)
			if !fsOk {
				t.Fatalf("expected *ByteStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertUInt8SliceEqual(t, fs.Data, []uint8{1, 10, 100, 255})
		})
	}
}

func TestBoolTensors(t *testing.T) {
	for _, filename := range makeFilenames("tensor_bool") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*BoolStorage)
			if !fsOk {
				t.Fatalf("expected *ByteStorage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertBoolSliceEqual(t, fs.Data, []bool{true, false, true, false})
		})
	}
}

func TestBFloat16Tensors(t *testing.T) { // Float
	for _, filename := range makeFilenames("tensor_bfloat16") {
		t.Run(filename, func(t *testing.T) {
			tensor := loadTensorFromFile(t, filename)
			assertCommonTensorFields(t, tensor)
			fs, fsOk := tensor.Source.(*BFloat16Storage)
			if !fsOk {
				t.Fatalf("expected *BFloat16Storage, got %#v", tensor.Source)
			}
			assertBaseStorageFields(t, fs.BaseStorage, 4, "cpu")
			assertFloat32SliceEqual(t, fs.Data, []float32{1.2, -3.4, 5.6, -7.8}, 0.013)
		})
	}
}

func loadTensorFromFile(t *testing.T, filename string) *Tensor {
	result, err := Load(path.Join("testdata", filename))
	if err != nil {
		t.Fatal(err)
	}
	tensor, tensorOk := result.(*Tensor)
	if !tensorOk {
		t.Fatalf("expected *Tensor, got %#v", result)
	}
	return tensor
}

func makeFilenames(prefix string) []string {
	filenames := make([]string, 0, 10)
	for _, proto := range [...]int{1, 2, 3, 4, 5} {
		for _, useZip := range [...]bool{false, true} {
			suffix := fmt.Sprintf("_proto%d", proto)
			if useZip {
				suffix += "_zip"
			}
			filenames = append(filenames, prefix+suffix+".pt")
		}
	}
	return filenames
}

func assertIntSliceEqual(t *testing.T, actual, expected []int) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		if actualVal != expected[index] {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertFloat32SliceEqual(t *testing.T, actual, expected []float32, eps float32) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		expectedVal := expected[index]
		if actualVal < expectedVal-eps || actualVal > expectedVal+eps {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertFloat64SliceEqual(t *testing.T, actual, expected []float64, eps float64) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		expectedVal := expected[index]
		if actualVal < expectedVal-eps || actualVal > expectedVal+eps {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertInt8SliceEqual(t *testing.T, actual, expected []int8) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		if actualVal != expected[index] {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertInt16SliceEqual(t *testing.T, actual, expected []int16) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		if actualVal != expected[index] {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertInt32SliceEqual(t *testing.T, actual, expected []int32) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		if actualVal != expected[index] {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertInt64SliceEqual(t *testing.T, actual, expected []int64) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		if actualVal != expected[index] {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertUInt8SliceEqual(t *testing.T, actual, expected []uint8) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		if actualVal != expected[index] {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertBoolSliceEqual(t *testing.T, actual, expected []bool) {
	if len(actual) != len(expected) {
		t.Errorf("expected %v, actual %v", expected, actual)
		return
	}
	for index, actualVal := range actual {
		if actualVal != expected[index] {
			t.Errorf("expected %v, actual %v", expected, actual)
			return
		}
	}
}

func assertCommonTensorFields(t *testing.T, tensor *Tensor) {
	assertIntSliceEqual(t, tensor.Size, []int{4})
	assertIntSliceEqual(t, tensor.Stride, []int{1})
	if tensor.StorageOffset != 0 {
		t.Errorf("expected StorageOffset 0, got %d", tensor.StorageOffset)
	}
	if tensor.RequiresGrad {
		t.Errorf("expected RequiresGrad false, got True")
	}
}

func assertBaseStorageFields(t *testing.T, bs BaseStorage, size int, location string) {
	if bs.Size != size {
		t.Errorf("expected storage Size %d, got %d", size, bs.Size)
	}
	if bs.Location != location {
		t.Errorf("expected storage Location %#v, got %#v", location, bs.Location)
	}
}
