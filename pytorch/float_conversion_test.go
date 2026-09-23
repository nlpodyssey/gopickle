package pytorch

import (
	"math"
	"testing"
)

func TestHalfFloatRoundTrip(t *testing.T) {
	t.Parallel()

	for i := 0; i <= math.MaxUint16; i++ {
		original := uint16(i)
		converted := FloatBits16to32(original)
		restored := FloatBits32to16(converted)
		if original != restored {
			t.Errorf("0b%016b was restored as 0b%016b", original, restored)
		}
	}
}

func TestFloatBits16to32(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name         string
		bits         uint16
		expected     float64
		wantNaN      bool
		wantNegative bool
	}{
		{name: "zero",
			bits: 0b0_00000_0000000000, expected: 0},
		{name: "negative zero",
			bits: 0b1_00000_0000000000, expected: 0, wantNegative: true},
		{name: "Inf",
			bits: 0b0_11111_0000000000, expected: math.Inf(1)},
		{name: "negative Inf",
			bits: 0b1_11111_0000000000, expected: math.Inf(-1), wantNegative: true},
		{name: "NaN",
			bits: 0b0_11111_0000000001, wantNaN: true},
		{name: "big NaN",
			bits: 0b0_11111_1111111111, wantNaN: true},
		{name: "negative NaN",
			bits: 0b1_11111_0000000001, wantNaN: true, wantNegative: true},
		{name: "negative big NaN",
			bits: 0b1_11111_1111111111, wantNaN: true, wantNegative: true},
		{name: "smallest normal",
			bits: 0b0_00001_0000000000, expected: 0x1p-14},
		{name: "smallest negative normal",
			bits: 0b1_00001_0000000000, expected: -0x1p-14, wantNegative: true},
		{name: "largest normal",
			bits: 0b0_11110_1111111111, expected: 0x1p15 * (1 + (1 - 0x1p-10))},
		{name: "largest negative normal",
			bits: 0b1_11110_1111111111, expected: -0x1p15 * (1 + (1 - 0x1p-10)), wantNegative: true},
		{name: "smallest subnormal",
			bits: 0b0_00000_0000000001, expected: 0x1p-14 * 0x1p-10},
		{name: "smallest negative subnormal",
			bits: 0b1_00000_0000000001, expected: -0x1p-14 * 0x1p-10, wantNegative: true},
		{name: "largest subnormal",
			bits: 0b0_00000_1111111111, expected: 0x1p-14 * (1 - 0x1p-10)},
		{name: "largest negative subnormal",
			bits: 0b1_00000_1111111111, expected: -0x1p-14 * (1 - 0x1p-10), wantNegative: true},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			converted := FloatBits16to32(tc.bits)
			actual := float64(math.Float32frombits(converted))
			if tc.wantNaN {
				if !math.IsNaN(actual) {
					t.Errorf("FloatBits16to32(0b%016b) = %f, want NaN with negative = %t",
						tc.bits, actual, tc.wantNegative)
				}
			} else {
				if actual != tc.expected {
					t.Errorf("FloatBits16to32(0b%016b) = %f, want %f with negative = %t",
						tc.bits, actual, tc.expected, tc.wantNegative)
				}
			}
			if math.Signbit(actual) != tc.wantNegative {
				t.Errorf("FloatBits16to32(0b%016b) negative = %t, want %t",
					tc.bits, math.Signbit(actual), tc.wantNegative)
			}
		})
	}
}

func TestFloatBits32to16(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		bits     uint32
		expected uint16
	}{
		{name: "zero",
			bits: 0b0_00000000_00000000000000000000000, expected: 0b0_00000_0000000000},
		{name: "negative zero",
			bits: 0b1_00000000_00000000000000000000000, expected: 0b1_00000_0000000000},
		{name: "too small subnormal",
			bits: 0b0_00000000_11111111111111111111111, expected: 0b0_00000_0000000000},
		{name: "too small",
			bits: 0b0_01100110_11111111111111111111111, expected: 0b0_00000_0000000000},
		{name: "too small negative",
			bits: 0b1_01100110_11111111111111111111111, expected: 0b1_00000_0000000000},
		{name: "smallest subnormal",
			bits: 0b0_01100111_00000000000000000000000, expected: 0b0_00000_0000000001},
		{name: "smallest negative subnormal",
			bits: 0b1_01100111_00000000000000000000000, expected: 0b1_00000_0000000001},
		{name: "smallest normal",
			bits: 0b0_01110001_00000000000000000000000, expected: 0b0_00001_0000000000},
		{name: "smallest negative normal",
			bits: 0b1_01110001_00000000000000000000000, expected: 0b1_00001_0000000000},
		{name: "one",
			bits: 0b0_01111111_00000000000000000000000, expected: 0b0_01111_0000000000},
		{name: "smallest after one",
			bits: 0b0_01111111_00000000010000000000000, expected: 0b0_01111_0000000001},
		{name: "rounding down",
			bits: 0b0_01111111_00000000001111111111111, expected: 0b0_01111_0000000000},
		{name: "largest",
			bits: 0b0_10001110_11111111110000000000000, expected: 0b0_11110_1111111111},
		{name: "largest negative",
			bits: 0b1_10001110_11111111110000000000000, expected: 0b1_11110_1111111111},
		{name: "positive overflow",
			bits: 0b0_10001111_00000000000000000000000, expected: 0b0_11111_0000000000},
		{name: "negative overflow",
			bits: 0b1_10001111_00000000000000000000000, expected: 0b1_11111_0000000000},
		{name: "NaN",
			bits: 0b0_11111111_01101101101010101010101, expected: 0b0_11111_0110110110},
		{name: "negative NaN",
			bits: 0b1_11111111_10000100001111011110111, expected: 0b1_11111_1000010000},
		{name: "MSB NaN",
			bits: 0b0_11111111_10000000000000000000000, expected: 0b0_11111_1000000000},
		{name: "negative MSB NaN",
			bits: 0b1_11111111_10000000000000000000000, expected: 0b1_11111_1000000000},
		{name: "LSB NaN",
			bits: 0b0_11111111_00000000000000000000001, expected: 0b0_11111_0000000001},
		{name: "negative LSB NaN",
			bits: 0b1_11111111_00000000000000000000001, expected: 0b1_11111_0000000001},
	} {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			converted := FloatBits32to16(tc.bits)
			if tc.expected != converted {
				t.Errorf("FloatBits32to16(0b%032b) = 0b%016b, want 0b%016b",
					tc.bits, converted, tc.expected)
			}
		})
	}

	t.Run("Go NaN", func(t *testing.T) {
		t.Parallel()

		converted := FloatBits32to16(math.Float32bits(float32(math.NaN())))
		if (converted&0b1_11111_0000000000) != 0b0_11111_0000000000 ||
			(converted&0b0_00000_1111111111) == 0 {
			t.Errorf("FloatBits32to16(math.NaN()) = 0b%016b, want NaN", converted)
		}
	})
}
