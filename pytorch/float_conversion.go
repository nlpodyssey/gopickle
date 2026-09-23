// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

// FloatBits16to32 converts the bits representation of a Half Float (16 bits)
// number to an IEEE 754 float representation (32 bits)
// From http://www.fox-toolkit.org/ftp/fasthalffloatconversion.pdf
func FloatBits16to32(u16 uint16) uint32 {
	return mantissaTable[offsetTable[u16>>10]+(uint32(u16)&0x3ff)] + exponentTable[u16>>10]
}

// FloatBits32to16 converts the bits representation of an IEEE 754 float
// representation (32 bits) to a half float (16 bits). It is a precise,
// bit-for-bit inverse of FloatBits16to32 for all valid half float values.
// Floats that cannot be represented as half floats are rounded using an
// unspecified method.
//
// The paper that FloatBits16to32 and FloatBits32to16 is based on contains an
// error: 32-bit NaNs with payloads beginning with 10 zero bits would become
// infinity if truncated to 16 bits according to the algorithm in the paper.
// We adopt the approach of truncating the least significant payload bits so
// that FloatBits32to16 precisely inverts FloatBits16to32, but we add a special
// case: 32-bit NaNs with payloads that would become 0 after truncation are
// instead changed to have a payload of 1, which preserves the NaN.
func FloatBits32to16(u32 uint32) uint16 {
	exp := (u32 >> 23) & 0x1ff
	mantissa32 := u32 & 0x007fffff
	mantissa16 := uint16(mantissa32 >> shiftTable[exp])
	out := baseTable[exp] + mantissa16
	if (exp&0xff) == 0xff && mantissa16 == 0 && mantissa32 != 0 {
		out |= 1 // NaN needs a payload
	}
	return out
}

// Tables for half -> float

var mantissaTable [2048]uint32
var exponentTable [64]uint32
var offsetTable [64]uint32

// Tables for float -> half

var baseTable [512]uint16
var shiftTable [512]uint8

func init() {
	initMantissaTable()
	initExponentTable()
	initOffsetTable()
	initBaseTable()
	initShiftTable()
}

func initMantissaTable() {
	mantissaTable[0] = 0
	for i := uint32(1); i < 1024; i++ {
		mantissaTable[i] = convertMantissa(i)
	}
	for i := uint32(1024); i < 2048; i++ {
		mantissaTable[i] = 0x38000000 + ((i - 1024) << 13)
	}
}

func initExponentTable() {
	exponentTable[0] = 0
	exponentTable[31] = 0x47800000
	exponentTable[32] = 0x80000000
	exponentTable[63] = 0xC7800000
	for i := uint32(1); i < 31; i++ {
		exponentTable[i] = i << 23
	}
	for i := uint32(33); i < 63; i++ {
		exponentTable[i] = 0x80000000 + (i-32)<<23
	}
}

func initOffsetTable() {
	offsetTable[0] = 0
	offsetTable[32] = 0
	for i := uint32(1); i < 32; i++ {
		offsetTable[i] = 1024
	}
	for i := uint32(33); i < 64; i++ {
		offsetTable[i] = 1024
	}
}

func convertMantissa(i uint32) uint32 {
	var m uint32 = i << 13  // zero pad mantissa bits
	var e uint32 = 0        // zero exponent
	for m&0x00800000 == 0 { // while not normalized
		e -= 0x00800000 // decrement exponent (1 << 23)
		m <<= 1         // shift mantissa
	}
	m &= ^uint32(0x00800000) // clear leading 1 bit
	e += 0x38800000          // adjust bias ((127-14)<<23)
	return m | e             // return combined number
}

func initBaseTable() {
	for i := int16(0); i < 256; i++ {
		e := i - 127
		switch {
		case e < -24: // Very small numbers map to zero
			baseTable[i|0x000] = 0x0000
			baseTable[i|0x100] = 0x8000
		case e < -14: // Small numbers map to denorms
			baseTable[i|0x000] = 0x0400 >> uint16(-e-14)
			baseTable[i|0x100] = (0x0400 >> uint16(-e-14)) | 0x8000
		case e <= 15: // Normal numbers just lose precision
			baseTable[i|0x000] = uint16((e + 15) << 10)
			baseTable[i|0x100] = uint16((e+15)<<10) | 0x8000
		case e < 128: // Large numbers map to Infinity
			baseTable[i|0x000] = 0x7c00
			baseTable[i|0x100] = 0xfc00
		default: // Infinity and NaN's stay Infinity and NaN's
			baseTable[i|0x000] = 0x7c00
			baseTable[i|0x100] = 0xfc00
		}
	}
}

func initShiftTable() {
	for i := int16(0); i < 256; i++ {
		e := i - 127
		switch {
		case e < -24:
			shiftTable[i|0x000] = 24
			shiftTable[i|0x100] = 24
		case e < -14:
			shiftTable[i|0x000] = uint8(-e) - 1
			shiftTable[i|0x100] = uint8(-e) - 1
		case e <= 15:
			shiftTable[i|0x000] = 13
			shiftTable[i|0x100] = 13
		case e < 128:
			shiftTable[i|0x000] = 24
			shiftTable[i|0x100] = 24
		default:
			shiftTable[i|0x000] = 13
			shiftTable[i|0x100] = 13
		}
	}
}
