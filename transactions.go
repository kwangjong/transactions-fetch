package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

const TIMESTAMPLAYOUT = "2006-01-02T15:04:05Z"

type Transaction struct {
	payer  string
	points int
	time   time.Time
}

type Balance struct {
	payer  string
	points int
}

func main() {
	// 0. check if arguments are valid
	if len(os.Args) != 3 {
		log.Fatal("Invalid arguments\nUsage:\n    transactions {points} {csv_path}")
	}

	point_to_spend, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("{points} should be integer\n Usage:\n    transactions {points} {csv_path}")
	}
	csv_path := os.Args[2]

	// 1. parse csv
	file, err := os.Open(csv_path)
	if err != nil {
		log.Fatal(err)
	}

	transactions := []Transaction{}

	parser := csv.NewReader(file)
	_, err = parser.Read() //skip first line
	if err != nil {
		log.Fatal(err)
	}
	for {
		line, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		points, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}

		time, err := time.Parse(TIMESTAMPLAYOUT, line[2])
		if err != nil {
			log.Fatal(err)
		}

		transactions = append(transactions, Transaction{
			payer:  line[0],
			points: points,
			time:   time,
		})
	}

	// 2. sort by timestamp
	sort.SliceStable(transactions, func(i, j int) bool {
		return transactions[i].time.Before(transactions[j].time)
	})

	// 3. process transactions and spend points
	balance := map[string]int{}

	for _, t := range transactions {
		if _, ok := balance[t.payer]; !ok {
			balance[t.payer] = 0
		}
		if point_to_spend > 0 {
			if point_to_spend-t.points > 0 {
				point_to_spend -= t.points
			} else {
				balance[t.payer] = t.points - point_to_spend
				point_to_spend = 0
			}
		} else {
			balance[t.payer] += t.points
		}
	}

	// 4. parse output
	jsonStr, err := json.MarshalIndent(balance, "", "\t")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(jsonStr))
	}
}
