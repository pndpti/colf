// Package colf writes and reads colf files. See docs/FORMAT.md.
package colf

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/pndpti/colf/schema"
)

type colInfo struct {
	offset int
	size   int
}

// Write encodes cols (each *schema.Int64Column or *schema.StringColumn,
// in the same order as sc.Fields) as a colf file.
func Write(w io.Writer, sc schema.Schema, cols []any) error {
	// TODO: magic + version, columns (record offset/size), footer, footer length, magic
	b := []byte("COLF")
	b = append(b, 1)
	var infos []colInfo
	numRows := 0

	for _, col := range cols {
		offset := len(b)
		switch c := col.(type) {
		case *schema.Int64Column:
			numRows = len(c.Values)
			for _, v := range c.Values {
				b = binary.LittleEndian.AppendUint64(b, uint64(v))
			}
		case *schema.StringColumn:
			numRows = len(c.Values)
			for _, s := range c.Values {
				b = binary.AppendUvarint(b, uint64(len(s)))
				b = append(b, s...)
			}

		}
		infos = append(infos, colInfo{offset: offset, size: len(b) - offset})
	}

	footerStart := len(b)
	b = binary.AppendUvarint(b, uint64(numRows))
	b = binary.AppendUvarint(b, uint64(len(cols)))
	for i, info := range infos {
		name := sc.Fields[i].Name
		typ := sc.Fields[i].Type
		b = binary.AppendUvarint(b, uint64(len(name)))
		b = append(b, name...)
		b = append(b, byte(typ))
		b = binary.AppendUvarint(b, uint64(info.offset))
		b = binary.AppendUvarint(b, uint64(info.size))
	}
	footerLength := len(b) - footerStart
	b = binary.LittleEndian.AppendUint32(b, uint32(footerLength))
	b = append(b, "COLF"...)

	_, err := w.Write(b)
	return err
}

// Read decodes a whole colf file.
func Read(b []byte) (schema.Schema, []any, error) {
	// TODO: check magic at both ends, read footer length, parse footer, decode each column
	return schema.Schema{}, nil, errors.New("not implemented")
}
