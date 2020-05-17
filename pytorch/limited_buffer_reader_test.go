// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

import (
	"bytes"
	"io"
	"testing"
)

func TestLimitedBufferReader(t *testing.T) {
	t.Run("empty input, data size 0", func(t *testing.T) {
		r := bytes.NewReader([]byte{})
		br := NewLimitedBufferReader(r, 0, 2, 3)
		assertNotHasNext(t, br)
		assertReadNextEof(t, br)
	})

	t.Run("empty input, data size > 0", func(t *testing.T) {
		r := bytes.NewReader([]byte{})
		br := NewLimitedBufferReader(r, 10, 2, 3)
		assertHasNext(t, br)
		assertReadNextEof(t, br)
	})

	t.Run("data size = scalar * buffer", func(t *testing.T) {
		input := []byte{11, 12, 21, 22, 31, 32, 41, 42, 51, 52, 61, 62}
		r := bytes.NewReader(input)
		br := NewLimitedBufferReader(r, 6, 2, 3)
		assertReadNextValue(t, br, []byte{11, 12})
		assertReadNextValue(t, br, []byte{21, 22})
		assertReadNextValue(t, br, []byte{31, 32})
		assertReadNextValue(t, br, []byte{41, 42})
		assertReadNextValue(t, br, []byte{51, 52})
		assertReadNextValue(t, br, []byte{61, 62})
		assertNotHasNext(t, br)
		assertReadNextEof(t, br)
	})

	t.Run("data size > scalar * buffer", func(t *testing.T) {
		input := []byte{11, 12, 21, 22, 31, 32, 41, 42, 51, 52, 61, 62, 71, 72}
		r := bytes.NewReader(input)
		br := NewLimitedBufferReader(r, 7, 2, 3)
		assertReadNextValue(t, br, []byte{11, 12})
		assertReadNextValue(t, br, []byte{21, 22})
		assertReadNextValue(t, br, []byte{31, 32})
		assertReadNextValue(t, br, []byte{41, 42})
		assertReadNextValue(t, br, []byte{51, 52})
		assertReadNextValue(t, br, []byte{61, 62})
		assertReadNextValue(t, br, []byte{71, 72})
		assertNotHasNext(t, br)
		assertReadNextEof(t, br)
	})

	t.Run("data size < scalar * buffer", func(t *testing.T) {
		input := []byte{11, 12, 21, 22, 31, 32, 41, 42, 51, 52}
		r := bytes.NewReader(input)
		br := NewLimitedBufferReader(r, 5, 2, 3)
		assertReadNextValue(t, br, []byte{11, 12})
		assertReadNextValue(t, br, []byte{21, 22})
		assertReadNextValue(t, br, []byte{31, 32})
		assertReadNextValue(t, br, []byte{41, 42})
		assertReadNextValue(t, br, []byte{51, 52})
		assertNotHasNext(t, br)
		assertReadNextEof(t, br)
	})

	t.Run("remaining data in buffer", func(t *testing.T) {
		input := []byte{11, 12, 21, 22, 31, 32, 41, 42, 51, 52, 90, 91, 92}
		r := bytes.NewReader(input)
		br := NewLimitedBufferReader(r, 5, 2, 3)
		assertReadNextValue(t, br, []byte{11, 12})
		assertReadNextValue(t, br, []byte{21, 22})
		assertReadNextValue(t, br, []byte{31, 32})
		assertReadNextValue(t, br, []byte{41, 42})
		assertReadNextValue(t, br, []byte{51, 52})
		assertNotHasNext(t, br)
		assertReadNextEof(t, br)

		rest := make([]byte, 3)
		n, err := r.Read(rest)
		if n != 3 || err != nil {
			t.Errorf("expected 3 bytes and no error, got %d and %#v", n, err)
		}
		assertByteSliceEqual(t, rest, []byte{90, 91, 92})
	})
}

func assertHasNext(t *testing.T, br *LimitedBufferReader) {
	if !br.HasNext() {
		t.Errorf("expected HasNext() true, but it is false")
	}
}

func assertNotHasNext(t *testing.T, br *LimitedBufferReader) {
	if br.HasNext() {
		t.Errorf("expected HasNext() false, but it is true")
	}
}

func assertReadNextValue(t *testing.T, br *LimitedBufferReader, val []byte) {
	assertHasNext(t, br)
	result, err := br.ReadNext()
	if err != nil {
		t.Errorf("expected nil error, actual %#v", err)
	}
	assertByteSliceEqual(t, result, val)
}

func assertByteSliceEqual(t *testing.T, actual, expected []byte) {
	if len(expected) != len(actual) {
		t.Errorf("expected %#v, actual %#v", expected, actual)
		return
	}
	for index, expected := range expected {
		actual := actual[index]
		if expected != actual {
			t.Errorf("expected %#v, actual %#v", expected, actual)
			return
		}
	}
}

func assertReadNextEof(t *testing.T, br *LimitedBufferReader) {
	result, err := br.ReadNext()
	if result != nil {
		t.Errorf("expected nil result, actual %#v", result)
	}
	if err != io.EOF {
		t.Errorf("expected EOF error, actual %#v", err)
	}
}
