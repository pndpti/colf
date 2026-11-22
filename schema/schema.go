package schema

type Type uint8

const (
	Int64 Type = iota
	String
)

type Field struct {
	Name string
	Type Type
}

type Schema struct {
	Fields []Field
}

type Int64Column struct {
	Values []int64
}

type StringColumn struct {
	Values []string
}
