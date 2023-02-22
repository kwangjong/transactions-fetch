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

type transaction struct {
	payer     string
	points    int
	timestamp time.Time
}

type Transaction []transaction

func (t Transaction) Len() int           { return len(t) }
func (t Transaction) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Transaction) Less(i, j int) bool { return t[i].timestamp.Before(t[j].timestamp) }

func read_transactions(filename string) *[]transaction {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Unable to open file")
	}
	defer file.Close()

	parser := csv.NewReader(file)
	parser.TrimLeadingSpace = true
	parser.FieldsPerRecord = 3

	_, err = parser.Read()
	if err != nil {
		log.Fatal("Unable to read CSV file")
	}

	var list_transactions []transaction

	for {
		line, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Unable to read CSV file")
		}

		payer := line[0]
		points, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal("Invalid points in CSV file")
		}

		timestamp, err := time.Parse(time.RFC3339, line[2])
		if err != nil {
			log.Fatal("Invalid timestamp in CSV file")
		}

		list_transactions = append(list_transactions, transaction{payer, points, timestamp})
	}

	return &list_transactions
}

func spend_points(points_to_spend int, list_transactions *[]transaction) map[string]int {
	sort.Sort(Transaction(*list_transactions))

	remaining_points := points_to_spend
	balance := map[string]int{}

	for _, t := range *list_transactions {
		if _, ok := balance[t.payer]; !ok {
			balance[t.payer] = 0
		}
		if remaining_points <= 0 {
			balance[t.payer] += t.points
			continue
		}

		if remaining_points-t.points > 0 {
			remaining_points -= t.points
		} else {
			balance[t.payer] = t.points - remaining_points
			remaining_points = 0
		}
	}

	if remaining_points > 0 {
		fmt.Errorf("Insufficient points to spend: requested=%d, remaining=%d", points_to_spend, remaining_points)
	}

	return balance
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: transactions <points> <filename>")
		os.Exit(1)
	}

	points_to_spend, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Invalid points argument")
	}
	filename := os.Args[2]

	list_transactions := read_transactions(filename)
	balance := spend_points(points_to_spend, list_transactions)

	jsonStr, err := json.MarshalIndent(balance, "", "\t")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(jsonStr))
	}
}
