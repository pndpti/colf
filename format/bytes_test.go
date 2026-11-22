package format

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

func TestBufReaderBytes(t *testing.T) {
	r := NewBufReader([]byte{1, 2, 3, 4, 5})

	if got := r.Bytes(2); !bytes.Equal(got, []byte{1, 2}) {
		t.Fatalf("first read = %v, want [1 2]", got)
	}
	if got := r.Bytes(3); !bytes.Equal(got, []byte{3, 4, 5}) {
		t.Fatalf("second read = %v, want [3 4 5]", got)
	}
	if r.Err() != nil {
		t.Fatalf("unexpected error: %v", r.Err())
	}
}

func TestReaderBufErrors(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		n    int
	}{
		{"past end", []byte{1, 2}, 3},
		{"empty input", nil, 1},
		{"negative", []byte{1, 2}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewBufReader(tt.data)
			if got := r.Bytes(tt.n); got != nil {
				t.Errorf("got %v, want nil", got)
			}
			if !errors.Is(r.Err(), io.ErrUnexpectedEOF) {
				t.Errorf("err = %v, want io.ErrUnexpectedEOF", r.Err())
			}
		})
	}
}

func TestBufReaderErrorSticks(t *testing.T) {
	r := NewBufReader([]byte{1, 2, 3})
	r.Bytes(10)
	if got := r.Bytes(1); got != nil {
		t.Errorf("read after error returned %v, want nil", got)
	}
	if r.Err() == nil {
		t.Error("error was cleared")
	}
}

func TestBufReaderUint32(t *testing.T) {
	buf := binary.LittleEndian.AppendUint32(nil, 0xABCDEFAD)
	r := NewBufReader(buf)
	if got := r.Uint32(); got != 0xABCDEFAD {
		t.Fatalf("got %x", got)
	}
	if r.Err() != nil {
		t.Fatal(r.Err())
	}
}
