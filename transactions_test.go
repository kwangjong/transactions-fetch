package main

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const (
	PATH_TEST_DIR      = "test_csv"
	PATH_VALID         = "transactions1.csv"
	PATH_REVERSE_ORDER = "transactions2.csv"
	PATH_INVALID       = "transactions3.csv"
)

var transactions_valid *[]transaction
var transactions_reverse_order *[]transaction
var transactions_invalid *[]transaction

// Test case 1: Spend some points from each player
func Test_spend_some_points(t *testing.T) {
	output, err := spend_points(5000, transactions_valid)
	expected := map[string]int{
		"DANNON":       1000,
		"MILLER COORS": 5300,
		"UNILEVER":     0,
	}

	if err != nil {
		t.Errorf("%v", err)
	}

	assert(t, output, expected)
}

// Test case 2: Spend all points
func Test_spend_all_points(t *testing.T) {
	output, err := spend_points(11300, transactions_valid)
	expected := map[string]int{
		"DANNON":       0,
		"MILLER COORS": 0,
		"UNILEVER":     0,
	}

	if err != nil {
		t.Error(err)
	}
	assert(t, output, expected)
}

// Test case 3: Spend points more than total balance of all payer
func Test_spend_more_points(t *testing.T) {
	_, err := spend_points(20000, transactions_valid)

	if !strings.Contains(err.Error(), "Insufficient balance") {
		t.Errorf("Wrong Error Thrown: Expected: Insufficient balance Error: %s", err.Error())
	}
}

// Test case 4: Spend all points, but transactions are in reverse chronological order
func Test_spend_reverse_order(t *testing.T) {
	output, err := spend_points(5000, transactions_reverse_order)
	expected := map[string]int{
		"DANNON":       500,
		"MILLER COORS": 3000,
		"UNILEVER":     500,
		"UW":           55000,
	}
	if err != nil {
		t.Error(err)
	}
	assert(t, output, expected)
}

// Test case 5: Spend all points, but transactions are invalid, i.e. negative transactions more than payer's balance
func Test_spend_invalid_transaction(t *testing.T) {
	_, err := spend_points(20000, transactions_valid)

	if strings.Contains(err.Error(), "Invalid transaction") {
		t.Errorf("Wrong Error Thrown: Expected: Insufficient balance Error: %s", err.Error())
	}
}

func assert(t *testing.T, output map[string]int, expected map[string]int) {
	if !reflect.DeepEqual(output, expected) {
		t.Errorf("AssertionError: Expected: %v Output: %v", expected, output)
	}
}

func setup() error {
	var err error
	transactions_valid, err = read_transactions(filepath.Join(PATH_TEST_DIR, PATH_VALID))
	if err != nil {
		return err
	}

	transactions_reverse_order, err = read_transactions(filepath.Join(PATH_TEST_DIR, PATH_REVERSE_ORDER))
	if err != nil {
		return err
	}

	transactions_invalid, err = read_transactions(filepath.Join(PATH_TEST_DIR, PATH_INVALID))
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	os.Exit(code)
}
