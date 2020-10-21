package deposits

import (
	"limits/database"
	"time"
)

const (
	period1W = "1w"
	period1D = "1d"
)

// DepositGetter holds logic to read deposits from the database.
type DepositGetter struct {
	db *database.Database
}

// NewDepositGetter creates a new instance of DepositGetter.
func NewDepositGetter(db *database.Database) *DepositGetter {
	return &DepositGetter{db}
}

// OneDay gets all deposits for a user that happened within a day before a certain date.
func (g *DepositGetter) OneDay(customerID string, before time.Time) []*database.Deposit {
	return g.get(customerID, period1D, before)
}

// OneWeek gets all deposits for a user that happened within a day before a certain date.
func (g *DepositGetter) OneWeek(customerID string, before time.Time) []*database.Deposit {
	return g.get(customerID, period1W, before)
}

// Exists checks if the deposit already exists in the database.
func (g *DepositGetter) Exists(customerID, depositID string) bool {
	_, ok := g.db.Customers[customerID]
	if !ok {
		return false
	}

	_, ok = g.db.Customers[customerID].Deposits[depositID]
	return ok
}

func (g *DepositGetter) get(customerID string, period string, before time.Time) []*database.Deposit {
	res := []*database.Deposit{}

	u, ok := g.db.Customers[customerID]
	if !ok {
		return res
	}

	for _, dep := range u.Deposits {
		if !dep.Accepted {
			continue
		}

		depYear, depWeek := dep.Time.ISOWeek()
		beforeYear, beforeWeek := before.ISOWeek()

		if period == period1D {
			if dep.Time.Day() == before.Day() && depYear == beforeYear && depWeek == beforeWeek {
				res = append(res, dep)
			}
		}

		if period == period1W {
			if depYear == beforeYear && depWeek == beforeWeek {
				res = append(res, dep)
			}
		}

	}

	return res
}
