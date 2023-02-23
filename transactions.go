package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

type transaction struct {
	payer     string
	points    int
	timestamp time.Time
}

type Transaction []transaction

func (t Transaction) Len() int           { return len(t) }
func (t Transaction) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Transaction) Less(i, j int) bool { return t[i].timestamp.Before(t[j].timestamp) }

func read_transactions(filename string) (*[]transaction, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("Unable to open file")
	}
	defer file.Close()

	parser := csv.NewReader(file)
	parser.TrimLeadingSpace = true
	parser.FieldsPerRecord = 3

	_, err = parser.Read()
	if err != nil {
		return nil, errors.New("Unable to read CSV file")
	}

	var list_transactions []transaction

	for {
		line, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("Unable to read CSV file")
		}

		payer := line[0]
		points, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, errors.New("Invalid points in CSV file")
		}

		timestamp, err := time.Parse(time.RFC3339, line[2])
		if err != nil {
			return nil, errors.New("Invalid timestamp in CSV file")
		}

		list_transactions = append(list_transactions, transaction{payer, points, timestamp})
	}

	return &list_transactions, nil
}

func spend_points(points_to_spend int, list_transactions *[]transaction) (map[string]int, error) {
	sort.Sort(Transaction(*list_transactions))

	remaining_points := points_to_spend
	balance_after_spending := map[string]int{}
	balance_before_spending := map[string]int{}

	for _, t := range *list_transactions {
		if _, ok := balance_after_spending[t.payer]; !ok {
			balance_after_spending[t.payer] = 0
			balance_before_spending[t.payer] = 0
		}

		balance_before_spending[t.payer] += t.points
		if balance_before_spending[t.payer] < 0 {
			return nil, errors.New("Invalid transaction: payer balance cannot be lower than zero")
		}

		if remaining_points <= 0 {
			balance_after_spending[t.payer] += t.points
		} else if remaining_points-t.points > 0 {
			remaining_points -= t.points
		} else {
			balance_after_spending[t.payer] = t.points - remaining_points
			remaining_points = 0
		}
	}

	if remaining_points > 0 {
		err := fmt.Errorf("Insufficient balance: requested=%d, remaining=%d", points_to_spend, remaining_points)
		return nil, err
	}

	return balance_after_spending, nil
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

	list_transactions, err := read_transactions(filename)
	if err != nil {
		log.Fatal(err)
	}
	balance, err := spend_points(points_to_spend, list_transactions)
	if err != nil {
		log.Fatal(err)
	}

	jsonStr, err := json.MarshalIndent(balance, "", "\t")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(string(jsonStr))
	}
}
