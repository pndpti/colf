package schema

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
