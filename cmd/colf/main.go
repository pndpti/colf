package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/pndpti/colf"
	"github.com/pndpti/colf/schema"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: colf write <in.csv> <out.colf> | colf read <file.colf> <column>")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "write":
		err = write(os.Args[2], os.Args[3])
	case "read":
		err = read(os.Args[2], os.Args[3])
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func write(in, out string) error {
	file, err := os.Open(in)
	if err != nil {
		return err
	}
	defer file.Close()
	r := csv.NewReader(file)
	_, err = r.Read()
	if err != nil {
		return err
	}

	sc := schema.Schema{Fields: []schema.Field{
		{Name: "order_id", Type: schema.Int64},
		{Name: "created_at", Type: schema.Int64},
		{Name: "customer", Type: schema.String},
		{Name: "city", Type: schema.String},
	}}

	idCol := &schema.Int64Column{}
	createdAtCol := &schema.Int64Column{}
	customerCol := &schema.StringColumn{}
	cityCol := &schema.StringColumn{}
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		id, err := strconv.ParseInt(rec[0], 10, 64)
		if err != nil {
			return err
		}
		idCol.Values = append(idCol.Values, id)

		createdAt, err := strconv.ParseInt(rec[1], 10, 64)
		if err != nil {
			return err
		}
		createdAtCol.Values = append(createdAtCol.Values, createdAt)
		customerCol.Values = append(customerCol.Values, rec[2])
		cityCol.Values = append(cityCol.Values, rec[3])
	}
	cols := []any{idCol, createdAtCol, customerCol, cityCol}

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := colf.Write(f, sc, cols); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func read(file, column string) error {
	return nil
}
