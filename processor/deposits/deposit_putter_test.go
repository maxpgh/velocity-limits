package deposits

import (
	"limits/database"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_NewDepositPutter(t *testing.T) {
	// Arrange
	db := database.New()

	// Act
	p := NewDepositPutter(db)

	// Assert
	require.NotNil(t, p)
}

func Test_DepositPutter_Put_CustomerNotExists(t *testing.T) {
	// Arrange
	db := database.New()
	p := NewDepositPutter(db)

	// Act
	dep := database.Deposit{
		ID:     "1",
		Amount: 1000,
		Time:   time.Now(),
	}
	p.Put("1", dep)

	// Assert
	require.NotEmpty(t, db.Customers)
	require.NotEmpty(t, db.Customers["1"].Deposits["1"])

	require.Equal(t, db.Customers["1"].Deposits["1"], &dep)
}

func Test_NewDepositPutter_Put_Success(t *testing.T) {
	// Arrange
	db := database.New()
	db.Customers["1"] = &database.Customer{
		ID:       "1",
		Deposits: make(map[string]*database.Deposit),
	}
	p := NewDepositPutter(db)

	// Act
	dep := database.Deposit{
		ID:     "1",
		Amount: 1000,
		Time:   time.Now(),
	}
	p.Put("1", dep)

	// Assert
	require.NotEmpty(t, db.Customers)
	require.NotEmpty(t, db.Customers["1"].Deposits["1"])
	require.Equal(t, db.Customers["1"].Deposits["1"], &dep)
}
