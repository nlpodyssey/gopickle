// Copyright 2023 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/binary"
	"fmt"
	"math"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

// Array unpickles array.array values as documented in:
//
//	https://docs.python.org/3/library/array.html
type Array struct{}

var _ Callable = (*Array)(nil)

func (Array) Call(args ...interface{}) (interface{}, error) {
	if got, want := len(args), 4; got != want {
		return nil, fmt.Errorf("invalid number of arguments (got=%d, want=%d)", got, want)
	}

	typ, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("invalid array type argument %T", args[1])
	}

	mi, ok := args[2].(int)
	if !ok {
		return nil, fmt.Errorf("invalid array mformat code type %T", args[2])
	}
	if mi >= len(arrayDescriptors) {
		return nil, fmt.Errorf("invalid array mformat value %d", mi)
	}
	descr := arrayDescriptors[mi]

	raw, ok := args[3].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid array payload type %T", args[3])
	}

	switch typ {
	case "b":
		vs := make([]int8, len(raw))
		for i := 0; i < len(raw); i++ {
			vs[i] = int8(raw[i])
		}
		return vs, nil

	case "B":
		return raw, nil

	case "u":
		vs := make([]rune, 0, utf8.RuneCount(raw))
		var enc encoding.Encoding
		switch descr.Size {
		case 4:
			order := unicode.BigEndian
			if descr.Order == binary.LittleEndian {
				order = unicode.LittleEndian
			}
			enc = unicode.UTF16(order, unicode.IgnoreBOM)
		case 8:
			order := utf32.BigEndian
			if descr.Order == binary.LittleEndian {
				order = utf32.LittleEndian
			}
			enc = utf32.UTF32(order, utf32.IgnoreBOM)
		default:
			return nil, fmt.Errorf("invalid machine description size (got=%d, want=4 or 8)", descr.Size)
		}
		dec := enc.NewDecoder()
		raw, err := dec.Bytes(raw)
		if err != nil {
			return nil, err
		}
		i := 0
	loop:
		for {
			r, sz := utf8.DecodeRune(raw[i:])
			switch r {
			case utf8.RuneError:
				if sz == 0 {
					break loop
				}
				return vs, fmt.Errorf("invalid rune")
			default:
				vs = append(vs, r)
				i += sz
			}
		}
		return vs, nil

	case "h":
		sz := descr.Size
		vs := make([]int16, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = int16(descr.Order.Uint16(raw[i:]))
		}
		return vs, nil

	case "H":
		sz := descr.Size
		vs := make([]uint16, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = descr.Order.Uint16(raw[i:])
		}
		return vs, nil

	case "i":
		sz := descr.Size
		vs := make([]int32, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = int32(descr.Order.Uint32(raw[i:]))
		}
		return vs, nil

	case "I":
		sz := descr.Size
		vs := make([]uint32, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = descr.Order.Uint32(raw[i:])
		}
		return vs, nil

	case "l":
		sz := descr.Size
		vs := make([]int64, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = int64(descr.Order.Uint64(raw[i:]))
		}
		return vs, nil

	case "L":
		sz := descr.Size
		vs := make([]uint64, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = descr.Order.Uint64(raw[i:])
		}
		return vs, nil

	case "q":
		sz := descr.Size
		vs := make([]int64, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = int64(descr.Order.Uint64(raw[i:]))
		}
		return vs, nil

	case "Q":
		sz := descr.Size
		vs := make([]uint64, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = descr.Order.Uint64(raw[i:])
		}
		return vs, nil

	case "f":
		sz := descr.Size
		vs := make([]float32, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = math.Float32frombits(descr.Order.Uint32(raw[i:]))
		}
		return vs, nil

	case "d":
		sz := descr.Size
		vs := make([]float64, len(raw)/sz)
		for i := 0; i < len(raw); i += sz {
			vs[i/sz] = math.Float64frombits(descr.Order.Uint64(raw[i:]))
		}
		return vs, nil

	default:
		return nil, fmt.Errorf("invalid array typecode '%s'", typ)
	}

	panic("impossible")
}

type arrayDescriptor struct {
	Size   int
	Signed bool
	Order  binary.ByteOrder
}

var (
	arrayDescriptors = []arrayDescriptor{
		0:  {Size: 1, Signed: false, Order: binary.LittleEndian}, // 0: UNSIGNED_INT8
		1:  {Size: 1, Signed: true, Order: binary.LittleEndian},  // 1: SIGNED_INT8
		2:  {Size: 2, Signed: false, Order: binary.LittleEndian}, // 2: UNSIGNED_INT16_LE
		3:  {Size: 2, Signed: false, Order: binary.BigEndian},    // 3: UNSIGNED_INT16_BE
		4:  {Size: 2, Signed: true, Order: binary.LittleEndian},  // 4: SIGNED_INT16_LE
		5:  {Size: 2, Signed: true, Order: binary.BigEndian},     // 5: SIGNED_INT16_BE
		6:  {Size: 4, Signed: false, Order: binary.LittleEndian}, // 6: UNSIGNED_INT32_LE
		7:  {Size: 4, Signed: false, Order: binary.BigEndian},    // 7: UNSIGNED_INT32_BE
		8:  {Size: 4, Signed: true, Order: binary.LittleEndian},  // 8: SIGNED_INT32_LE
		9:  {Size: 4, Signed: true, Order: binary.BigEndian},     // 9: SIGNED_INT32_BE
		10: {Size: 8, Signed: false, Order: binary.LittleEndian}, // 10: UNSIGNED_INT64_LE
		11: {Size: 8, Signed: false, Order: binary.BigEndian},    // 11: UNSIGNED_INT64_BE
		12: {Size: 8, Signed: true, Order: binary.LittleEndian},  // 12: SIGNED_INT64_LE
		13: {Size: 8, Signed: true, Order: binary.BigEndian},     // 13: SIGNED_INT64_BE
		14: {Size: 4, Signed: false, Order: binary.LittleEndian}, // 14: IEEE_754_FLOAT_LE
		15: {Size: 4, Signed: false, Order: binary.BigEndian},    // 15: IEEE_754_FLOAT_BE
		16: {Size: 8, Signed: false, Order: binary.LittleEndian}, // 16: IEEE_754_DOUBLE_LE
		17: {Size: 8, Signed: false, Order: binary.BigEndian},    // 17: IEEE_754_DOUBLE_BE
		18: {Size: 4, Signed: false, Order: binary.LittleEndian}, // 18: UTF16_LE
		19: {Size: 4, Signed: false, Order: binary.BigEndian},    // 19: UTF16_BE
		20: {Size: 8, Signed: false, Order: binary.LittleEndian}, // 20: UTF32_LE
		21: {Size: 8, Signed: false, Order: binary.BigEndian},    // 21: UTF32_BE
	}
)
