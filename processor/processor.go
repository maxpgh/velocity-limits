package processor

import (
	"encoding/json"
	"limits/database"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	limitDay   = 500000
	limitWeek  = 2000000
	limitTimes = 3
)

// Processor holds all business logic.
type Processor struct {
	db            *database.Database
	depositGetter DepositGetter
	DepositPutter DepositPutter
}

// New creates a new instance of Processor.
func New(
	db *database.Database,
	depositGetter DepositGetter,
	depositPutter DepositPutter,
) *Processor {
	return &Processor{
		db,
		depositGetter,
		depositPutter,
	}
}

// Input represents a line from the input file.
type Input struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	LoadAmount string    `json:"load_amount"`
	Time       time.Time `json:"time"`
}

// Output represents an line that goes to output file.
type Output struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted   bool   `json:"accepted"`
}

// DepositGetter abstracts database calls to get past deposits.
type DepositGetter interface {
	OneDay(customerID string, before time.Time) []*database.Deposit
	OneWeek(customerID string, before time.Time) []*database.Deposit
	Exists(customerID, depositID string) bool
}

// DepositPutter abstracts database calls to save incoming deposits.
type DepositPutter interface {
	Put(customer_id string, deposit database.Deposit)
}

// Process determines if an incoming deposit should be accepted and returns a JSON string.
func (p *Processor) Process(input []byte) ([]byte, error) {
	// parse the payload
	var inp Input

	err := json.Unmarshal(input, &inp)
	if err != nil {
		return nil, errors.Wrap(err, "limits: Processor.Process json.Unmarshal error")
	}

	if p.depositGetter.Exists(inp.CustomerID, inp.ID) {
		return nil, nil
	}

	amount, err := p.parseAmount(inp.LoadAmount)
	if err != nil {
		return nil, errors.Wrap(err, "limits: Processor.parseAmount error")
	}

	incomingDep := database.Deposit{
		ID:     inp.ID,
		Amount: amount,
		Time:   inp.Time,
	}

	// check if a user passed the depoists per day threshold
	deps := p.depositGetter.OneDay(inp.CustomerID, inp.Time)

	if len(deps) == limitTimes {
		incomingDep.Accepted = false
		p.DepositPutter.Put(inp.CustomerID, incomingDep)

		out := Output{
			ID:         inp.ID,
			CustomerID: inp.CustomerID,
			Accepted:   false,
		}

		res, err := json.Marshal(&out)
		if err != nil {
			return nil, errors.Wrap(err, "limits: Processor.Process json.Marshal error")
		}

		return res, nil
	}

	// check if a user deposited more than a threshold per day
	if (p.calculateAmount(deps) + amount) >= limitDay {
		incomingDep.Accepted = false
		p.DepositPutter.Put(inp.CustomerID, incomingDep)

		out := Output{
			ID:         inp.ID,
			CustomerID: inp.CustomerID,
			Accepted:   false,
		}

		res, err := json.Marshal(&out)
		if err != nil {
			return nil, errors.Wrap(err, "limits: Processor.Process json.Marshal error")
		}

		return res, nil
	}

	// check if a user deposited more than a threshold per week
	deps = p.depositGetter.OneWeek(inp.CustomerID, inp.Time)

	if (p.calculateAmount(deps) + amount) >= limitWeek {
		incomingDep.Accepted = false
		p.DepositPutter.Put(inp.CustomerID, incomingDep)

		out := Output{
			ID:         inp.ID,
			CustomerID: inp.CustomerID,
			Accepted:   false,
		}

		res, err := json.Marshal(&out)
		if err != nil {
			return nil, errors.Wrap(err, "limits: Processor.Process json.Marshal error")
		}

		return res, nil
	}

	// accepting the incoming deposit
	incomingDep.Accepted = true
	p.DepositPutter.Put(inp.CustomerID, incomingDep)

	out := Output{
		ID:         inp.ID,
		CustomerID: inp.CustomerID,
		Accepted:   true,
	}

	res, err := json.Marshal(&out)
	if err != nil {
		return nil, errors.Wrap(err, "limits: Processor.Process json.Marshal error")
	}

	return res, nil

}

func (p *Processor) calculateAmount(deps []*database.Deposit) int {
	amount := 0

	for _, dep := range deps {
		amount += dep.Amount
	}

	return amount
}

func (p *Processor) parseAmount(amount string) (int, error) {
	tmp := strings.TrimLeft(amount, "$")

	res, err := strconv.ParseFloat(tmp, 64)
	if err != nil {
		return 0, err
	}

	return int(res * 100), nil
}
