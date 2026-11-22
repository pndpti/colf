package colf

import (
	"bytes"
	"testing"

	"github.com/pndpti/colf/schema"
)

// TestWrite checks the exact bytes for id=[1], name=["a"].
// The layout follows docs/FORMAT.md.
func TestWrite(t *testing.T) {
	sc := schema.Schema{Fields: []schema.Field{
		{Name: "id", Type: schema.Int64},
		{Name: "name", Type: schema.String},
	}}
	cols := []any{
		&schema.Int64Column{Values: []int64{1}},
		&schema.StringColumn{Values: []string{"a"}},
	}

	want := bytes.Join([][]byte{
		[]byte("COLF"), {0x01}, // magic, version
		{0x01, 0, 0, 0, 0, 0, 0, 0},                  // id: 1
		{0x01, 'a'},                                  // name: len 1, "a"
		{0x01, 0x02},                                 // numRows=1, numCols=2
		{0x02, 'i', 'd', 0x00, 0x05, 0x08},           // "id", int64, offset 5, size 8
		{0x04, 'n', 'a', 'm', 'e', 0x01, 0x0D, 0x02}, // "name", string, offset 13, size 2
		{0x10, 0x00, 0x00, 0x00},                     // footer length = 16
		[]byte("COLF"),
	}, nil)

	var buf bytes.Buffer
	if err := Write(&buf, sc, cols); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if got := buf.Bytes(); !bytes.Equal(got, want) {
		t.Errorf("bytes mismatch\n got: % x\nwant: % x", got, want)
	}
}
