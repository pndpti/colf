// Package colf writes and reads colf files. See docs/FORMAT.md.
package colf

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/pndpti/colf/format"
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
	var sc schema.Schema
	var infos []colInfo

	// TODO: check magic at both ends, read footer length, parse footer, decode each column
	if len(b) < 13 || string(b[:4]) != "COLF" || string(b[len(b)-4:]) != "COLF" {
		return sc, nil, errors.New("not a colf file")
	}

	footerLength := int(binary.LittleEndian.Uint32(b[len(b)-8 : len(b)-4]))
	footerStart := len(b) - 8 - footerLength
	footer := b[footerStart : footerStart+footerLength]

	r := format.NewBufReader(footer)
	numRows := int(r.Uvarint())
	numCols := int(r.Uvarint())

	for range numCols {
		nameLength := int(r.Uvarint())
		name := string(r.Bytes(nameLength))
		typ := schema.Type(r.Bytes(1)[0])
		offset := int(r.Uvarint())
		size := int(r.Uvarint())

		sc.Fields = append(sc.Fields, schema.Field{Name: name, Type: typ})
		infos = append(infos, colInfo{offset: offset, size: size})
	}
	if r.Err() != nil {
		return sc, nil, r.Err()
	}

	var cols []any
	for i, info := range infos {
		cr := format.NewBufReader(b[info.offset : info.offset+info.size])

		switch sc.Fields[i].Type {
		case schema.Int64:
			c := &schema.Int64Column{}
			for range numRows {
				c.Values = append(c.Values, int64(cr.Uint64()))
			}
			cols = append(cols, c)
		case schema.String:
			c := &schema.StringColumn{}
			for range numRows {
				n := int(cr.Uvarint())
				s := string(cr.Bytes(n))
				c.Values = append(c.Values, s)
			}
			cols = append(cols, c)
		}
		if cr.Err() != nil {
			return sc, nil, cr.Err()
		}
	}

	return sc, cols, nil
}
