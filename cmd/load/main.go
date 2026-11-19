package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/pndpti/colf/schema"
)

func main() {
	in := flag.String("in", "", "input csv")
	schemaFlag := flag.String("schema", "", `e.g. "order_id:int64,city:string"`)
	flag.Parse()

	sc, err := schema.Parse(*schemaFlag)
	if err != nil {
		log.Fatal(err)
	}

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

	file, err := os.Open(*in)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { err = errors.Join(err, file.Close()) }()

	r := csv.NewReader(file)
	r.Read()

	rows := 0
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
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
				c.Valid = append(c.Valid, s != "")
			}
		}
		rows++
	}

	fmt.Println("loaded", rows, "rows")

	for i := 0; i < 5 && i < rows; i++ {
		for _, col := range cols {
			switch c := col.(type) {
			case *schema.Int64Column:
				fmt.Print(c.Values[i], " ")
			case *schema.TimestampColumn:
				fmt.Print(c.Values[i], " ")
			case *schema.Float64Column:
				fmt.Print(c.Values[i], " ")
			case *schema.BoolColumn:
				fmt.Print(c.Values[i], " ")
			case *schema.StringColumn:
				fmt.Print(c.Values[i], " ")
			}
		}
		fmt.Println()
	}
}
