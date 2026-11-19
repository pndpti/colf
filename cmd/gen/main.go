package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	rows := flag.Int("rows", 1_000_000, "number of rows to generate")
	out := flag.String("out", "sales.csv", "output file name")
	flag.Parse()

	rng := rand.New(rand.NewSource(42))

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
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	fmt.Fprintln(w, "order_id,created_at,customer,city,amount,is_returned")

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

		amount := ""
		if rng.Float64() < 0.05 {
			amount = strconv.FormatFloat(rng.Float64()*1000, 'f', 2, 64)
		}

		isReturned := rng.Float64() < 0.03

		fmt.Fprintf(w, "%d,%d,%s,%s,%s,%t\n",
			orderID, ts.Unix(), customer, city, amount, isReturned)
	}
}
