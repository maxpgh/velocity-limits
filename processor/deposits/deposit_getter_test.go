package deposits

import (
	"limits/database"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_NewDepositGetter(t *testing.T) {
	// Arrange
	db := database.New()

	// Act
	g := NewDepositGetter(db)

	// Assert
	require.NotNil(t, g)
}

func Test_DepositGetter_NoCustomers_EmptyResult(t *testing.T) {
	// Arrange
	db := database.New()
	g := NewDepositGetter(db)

	// Act
	res := g.OneWeek("1", time.Now())

	// Assert
	require.Empty(t, res)
}

func Test_DepositGetter_Exists(t *testing.T) {
	// Arrange
	db := &database.Database{
		Customers: map[string]*database.Customer{
			"1": {
				ID: "1",
				Deposits: map[string]*database.Deposit{
					"2": {
						ID:     "2",
						Amount: 1000,
						Time:   time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().UTC().Location()),
					},
				},
			},
		},
	}
	g := NewDepositGetter(db)

	// Act
	res := g.Exists("1", "2")

	// Assert
	require.True(t, res)
}

func Test_DepositGetter_OneDay_Success(t *testing.T) {
	// Arrange
	db := database.New()
	db.Customers["1"] = &database.Customer{
		ID: "1",
		Deposits: map[string]*database.Deposit{
			"1": {
				ID:       "1",
				Amount:   100,
				Time:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().UTC().Location()),
				Accepted: true,
			},
			"2": {
				ID:       "2",
				Amount:   100,
				Time:     time.Date(2000, 1, 1, 11, 15, 2, 0, time.Now().UTC().Location()),
				Accepted: true,
			},
		},
	}

	g := NewDepositGetter(db)

	// Act
	res := g.OneDay("1", time.Date(2000, 1, 1, 11, 15, 2, 0, time.Now().UTC().Location()))

	// Assert
	require.NotEmpty(t, res)
	require.Equal(t, len(res), 2)
	require.Equal(t, "1", res[0].ID)
	require.Equal(t, "2", res[1].ID)
}

func Test_DepositGetter_OneWeek_Success(t *testing.T) {
	// Arrange

	db := database.New()
	db.Customers["1"] = &database.Customer{
		ID: "1",
		Deposits: map[string]*database.Deposit{
			"1": {
				ID:       "1",
				Amount:   100,
				Time:     time.Date(2000, 1, 1, 0, 0, 0, 0, time.Now().UTC().Location()),
				Accepted: true,
			},
			"2": {
				ID:       "2",
				Amount:   100,
				Time:     time.Date(2000, 1, 1, 11, 15, 2, 0, time.Now().UTC().Location()),
				Accepted: true,
			},
			"3": {
				ID:       "3",
				Amount:   100,
				Time:     time.Date(2000, 1, 2, 19, 58, 46, 0, time.Now().UTC().Location()),
				Accepted: true,
			},
		},
	}

	g := NewDepositGetter(db)

	// Act
	res := g.OneWeek("1", time.Date(2000, 1, 2, 19, 58, 46, 0, time.Now().UTC().Location()))

	// Assert
	require.NotEmpty(t, res)
	require.Equal(t, len(res), 3)
}
