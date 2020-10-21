package deposits

import "limits/database"

// DepositPutter holds logic to write a deposit into database.
type DepositPutter struct {
	db *database.Database
}

// NewDepositPutter creates a new instance of DepositPutter.
func NewDepositPutter(db *database.Database) *DepositPutter {
	return &DepositPutter{db}
}

// Put records a new deposit into the database.
func (p *DepositPutter) Put(customer_id string, deposit database.Deposit) {
	u, ok := p.db.Customers[customer_id]

	if !ok {
		deps := map[string]*database.Deposit{
			deposit.ID: &deposit,
		}

		p.db.Customers[customer_id] = &database.Customer{
			ID:       customer_id,
			Deposits: deps,
		}

		return
	}

	u.Deposits[deposit.ID] = &deposit
}
