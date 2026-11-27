package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	rows := flag.Int("rows", 1_000_000, "number of rows to generate")
	out := flag.String("out", "sales.csv", "output file name")
	flag.Parse()

	rng := rand.New(rand.NewSource(42)) //nolint:gosec

	cities := []string{
		"Mumbai", "Delhi", "Bengaluru", "Hyderabad", "Chennai",
		"Kolkata", "Pune", "Ahmedabad", "Jaipur", "Surat",
		"Lucknow", "Kanpur", "Nagpur", "Indore", "Bhopal",
		"Nellore", "Vizag", "Kochi", "Patna", "Coimbatore",
	}

	customers := make([]string, 50000)
	for i := range customers {
		customers[i] = "customer_" + strconv.Itoa(i)
	}

	file, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer func() { err = errors.Join(err, file.Close()) }()

	w := bufio.NewWriter(file)
	defer func() { err = errors.Join(err, w.Flush()) }()

	_, err = fmt.Fprintln(w, "order_id,created_at,customer,city")
	if err != nil {
		log.Fatal(err)
	}

	orderID := int64(1000)
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	city := cities[0]

	for i := 0; i < *rows; i++ {
		orderID += int64(1 + rng.Intn(3))

		ts = ts.Add(time.Duration(1+rng.Intn(5)) * time.Second)

		customer := customers[rng.Intn(len(customers))]

		if rng.Float64() < 0.10 {
			city = cities[rng.Intn(len(cities))]
		}

		_, err = fmt.Fprintf(w, "%d,%d,%s,%s\n",
			orderID, ts.Unix(), customer, city)
		if err != nil {
			log.Fatal(err)
		}
	}
}
