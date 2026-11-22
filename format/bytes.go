package format

import (
	"encoding/binary"
	"errors"
	"io"
)

type BufReader struct {
	buf []byte
	err error
}

func NewBufReader(b []byte) *BufReader {
	return &BufReader{buf: b}
}

func (r *BufReader) Err() error {
	return r.err
}

func (r *BufReader) Bytes(n int) []byte {
	if r.err != nil {
		return nil
	}
	if n < 0 || n > len(r.buf) {
		r.err = io.ErrUnexpectedEOF
		return nil
	}
	b := r.buf[:n]
	r.buf = r.buf[n:]
	return b
}

func (r *BufReader) Uint32() uint32 {
	b := r.Bytes(4)
	if b == nil {
		return 0
	}
	return binary.LittleEndian.Uint32(b)
}

func (r *BufReader) Uvarint() uint64 {
	if r.err != nil {
		return 0
	}
	v, n := binary.Uvarint(r.buf)
	if n <= 0 {
		r.err = errors.New("format: bad varint")
		return 0
	}
	r.buf = r.buf[n:]
	return v
}
