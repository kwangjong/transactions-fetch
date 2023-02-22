package main

import "testing"

// Test case 1: Spend some points from each player
func Test_spend_some_points(t *testing.T)

// Test case 2: Spend all points
func Test_spend_all_points(t *testing.T)

// Test case 3: Spend points more than total balance of all payer
func Test_spend_more_points(t *testing.T)

// Test case 4: Spend all points, but transactions are in reverse chronological order
func Test_spend_reverse_order(t *testing.T)

// Test case 5: Spend all points, but transactions are invalid, i.e. negative transactions more than payer's balance
func Test_spend_invalid_transaction(t *testing.T)
