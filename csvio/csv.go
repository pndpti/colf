package csvio

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/pndpti/colf/schema"
)

func Load(in io.Reader, sc schema.Schema) ([]any, error) {
	var cols []any
	for _, f := range sc.Fields {
		switch f.Type {
		case schema.Int64:
			cols = append(cols, &schema.Int64Column{})
		case schema.Float64:
			cols = append(cols, &schema.Float64Column{})
		case schema.String:
			cols = append(cols, &schema.StringColumn{})
		case schema.Bool:
			cols = append(cols, &schema.BoolColumn{})
		case schema.Timestamp:
			cols = append(cols, &schema.TimestampColumn{})
		}
	}

	r := csv.NewReader(in)
	_, err := r.Read()
	if err != nil {
		return nil, err
	}

	rows := 0
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for i, col := range cols {
			s := rec[i]
			switch c := col.(type) {
			case *schema.Int64Column:
				v, _ := strconv.ParseInt(s, 10, 64)
				c.Values = append(c.Values, v)
				c.Valid = append(c.Valid, s != "")
			case *schema.TimestampColumn:
				v, _ := strconv.ParseInt(s, 10, 64)
				c.Values = append(c.Values, v)
				c.Valid = append(c.Valid, s != "")
			case *schema.Float64Column:
				v, _ := strconv.ParseFloat(s, 64)
				c.Values = append(c.Values, v)
				c.Valid = append(c.Valid, s != "")
			case *schema.BoolColumn:
				v, _ := strconv.ParseBool(s)
				c.Values = append(c.Values, v)
				c.Valid = append(c.Valid, s != "")
			case *schema.StringColumn:
				c.Values = append(c.Values, s)
				c.Valid = append(c.Valid, true)
			}
		}
		rows++
	}
	return cols, nil
}
