package schema

import (
	"fmt"
	"strings"
)

type Type uint8

const (
	Int64 Type = iota
	Float64
	String
	Bool
	Timestamp
)

type Field struct {
	Name     string
	Type     Type
	Nullable bool
}

type Schema struct {
	Fields []Field
}

type Int64Column struct {
	Values []int64
	Valid  []bool
}

type Float64Column struct {
	Values []float64
	Valid  []bool
}

type StringColumn struct {
	Values []string
	Valid  []bool
}

type BoolColumn struct {
	Values []bool
	Valid  []bool
}

type TimestampColumn struct {
	Values []int64
	Valid  []bool
}

func Parse(s string) (Schema, error) {
	var sc Schema
	for _, part := range strings.Split(s, ",") {
		name, typ, ok := strings.Cut(strings.TrimSpace(part), ":")
		if !ok {
			return sc, fmt.Errorf("bad field %q, want name:type", part)
		}
		var t Type
		switch typ {
		case "int64":
			t = Int64
		case "float64":
			t = Float64
		case "string":
			t = String
		case "bool":
			t = Bool
		case "timestamp":
			t = Timestamp
		default:
			return sc, fmt.Errorf("unknown type %q", typ)
		}
		sc.Fields = append(sc.Fields, Field{Name: name, Type: t, Nullable: true})
	}
	return sc, nil
}
