package database

import "time"

// Database holds all the data.
type Database struct {
	Customers map[string]*Customer
}

// Customer represents a customer.
type Customer struct {
	ID       string
	Deposits map[string]*Deposit
}

// Deposit represents a "load", incoming money from a customer.
type Deposit struct {
	ID       string
	Amount   int
	Time     time.Time
	Accepted bool
}

// New creates a new Database.
func New() *Database {
	return &Database{
		Customers: make(map[string]*Customer),
	}
}
