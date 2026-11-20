package format

import "io"

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
