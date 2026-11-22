package colf

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/pndpti/colf/schema"
)

func TestRoundTrip(t *testing.T) {
	sc := schema.Schema{Fields: []schema.Field{
		{Name: "id", Type: schema.Int64},
		{Name: "name", Type: schema.String},
	}}
	cols := []any{
		&schema.Int64Column{Values: []int64{1, 2, 3}},
		&schema.StringColumn{Values: []string{"a", "bb", "ccc"}},
	}

	var buf bytes.Buffer
	if err := Write(&buf, sc, cols); err != nil {
		t.Fatal(err)
	}

	gotSc, gotCols, err := Read(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotSc, sc) {
		t.Errorf("schema = %+v, want %+v", gotSc, sc)
	}
	if !reflect.DeepEqual(gotCols, cols) {
		t.Errorf("cols = %+v, want %+v", gotCols, cols)
	}
}
