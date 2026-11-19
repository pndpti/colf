package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Column struct {
	Name    string
	Type    string
	Ints    []int64
	Floats  []float64
	Strings []string
	Bools   []bool
	Valid   []bool
}

func main() {
	in := flag.String("in", "", "input csv")
	schema := flag.String("schema", "", `e.g. "order_id:int64,city:string"`)
	flag.Parse()

	var cols []*Column
	for _, part := range strings.Split(*schema, ",") {
		name, typ, ok := strings.Cut(part, ":")
		if !ok {
			log.Fatalf("bad schema field %q", part)
		}
		cols = append(cols, &Column{Name: name, Type: typ})
	}

	f, err := os.Open(*in)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Read()

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for i, c := range cols {
			s := rec[i]
			c.Valid = append(c.Valid, s != "")

			switch c.Type {
			case "int64":
				v, _ := strconv.ParseInt(s, 10, 64)
				c.Ints = append(c.Ints, v)
			case "float64":
				v, _ := strconv.ParseFloat(s, 64)
				c.Floats = append(c.Floats, v)
			case "bool":
				v, _ := strconv.ParseBool(s)
				c.Bools = append(c.Bools, v)
			case "string":
				c.Strings = append(c.Strings, s)
			}
		}
	}

	fmt.Println("loaded", len(cols[0].Valid), "rows")

	for i := 0; i < 5; i++ {
		for _, c := range cols {
			switch c.Type {
			case "int64":
				fmt.Print(c.Ints[i], " ")
			case "float64":
				fmt.Print(c.Floats[i], " ")
			case "bool":
				fmt.Print(c.Bools[i], " ")
			case "string":
				fmt.Print(c.Strings[i], " ")
			}
		}
		fmt.Println()
	}

	for i := 0; i < 5; i++ {
		for _, c := range cols {
			fmt.Println(c.Strings[i], " ")
		}
	}
}
