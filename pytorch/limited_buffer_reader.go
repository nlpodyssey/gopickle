// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pytorch

import "io"

type LimitedBufferReader struct {
	r              io.Reader
	scalarSize     int
	remainingBytes int
	buf            []byte
	bufIndex       int
}

func NewLimitedBufferReader(
	r io.Reader,
	dataSize, scalarSize, bufferSize int,
) *LimitedBufferReader {
	size := bufferSize * scalarSize
	return &LimitedBufferReader{
		r:              r,
		scalarSize:     scalarSize,
		remainingBytes: scalarSize * dataSize,
		buf:            make([]byte, size),
		bufIndex:       size,
	}
}

func (br *LimitedBufferReader) HasNext() bool {
	return br.remainingBytes != 0
}

func (br *LimitedBufferReader) ReadNext() ([]byte, error) {
	if br.remainingBytes == 0 {
		return nil, io.EOF
	}
	if br.bufIndex == len(br.buf) {
		br.bufIndex = 0
		if br.remainingBytes < len(br.buf) {
			br.buf = br.buf[0:br.remainingBytes]
		}
		_, err := br.r.Read(br.buf)
		if err != nil {
			return nil, err
		}
	}
	result := br.buf[br.bufIndex : br.bufIndex+br.scalarSize]
	br.bufIndex += br.scalarSize
	br.remainingBytes -= br.scalarSize
	return result, nil
}
