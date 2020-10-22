package processor

import (
	"limits/database"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	// Arrange
	getter := mockDepositGetter{}
	putter := mockDepositPutter{}

	// Act
	p := New(&getter, &putter)

	// Assert
	require.NotNil(t, p)
}

func Test_Process_InvalidInput_Error(t *testing.T) {
	// Arrange
	getter := mockDepositGetter{}
	putter := mockDepositPutter{}
	p := New(&getter, &putter)

	// Act
	res, err := p.Process([]byte("x"))

	// Assert
	require.Error(t, err)
	require.Equal(t, "limits: Processor.Process json.Unmarshal error: invalid character 'x' looking for beginning of value", err.Error())
	require.Nil(t, res)
}

func Test_Process_DepositExists(t *testing.T) {
	// Arrange
	getter := mockDepositGetter{
		ResExists: true,
	}
	putter := mockDepositPutter{}
	p := New(&getter, &putter)

	// Act
	res, err := p.Process([]byte("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3318.47\",\"time\":\"2000-01-01T00:00:00Z\"}"))

	// Assert
	require.NoError(t, err)
	require.Nil(t, res)
}

func Test_Process_DepositMoreThan3TimesIn1Day_Denied(t *testing.T) {
	// Arrange
	getter := mockDepositGetter{
		Res1Day: []*database.Deposit{
			{
				ID:     "1",
				Amount: 1000,
				Time:   time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:     "2",
				Amount: 1000,
				Time:   time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:     "3",
				Amount: 1000,
				Time:   time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	putter := mockDepositPutter{}
	p := New(&getter, &putter)

	// Act
	res, err := p.Process([]byte("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3318.47\",\"time\":\"2000-01-01T00:00:00Z\"}"))

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, getter.OneDayCalled)
	require.Equal(t, "{\"id\":\"15887\",\"customer_id\":\"528\",\"accepted\":false}", string(res))
}

func Test_Process_InvalidInputAmount_Error(t *testing.T) {
	// Arrange
	getter := mockDepositGetter{}
	putter := mockDepositPutter{}
	p := New(&getter, &putter)

	// Act
	res, err := p.Process([]byte("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"abc\",\"time\":\"2000-01-01T00:00:00Z\"}"))

	// Assert
	require.Error(t, err)
	require.Equal(t, "limits: Processor.parseAmount error: strconv.ParseFloat: parsing \"abc\": invalid syntax", err.Error())
	require.Nil(t, res)
}

func Test_Process_DepositMoreThanLimitPer1Day_Denied(t *testing.T) {
	// Arrange
	getter := mockDepositGetter{
		Res1Day: []*database.Deposit{
			{
				ID:     "1",
				Amount: 400000,
				Time:   time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	putter := mockDepositPutter{}
	p := New(&getter, &putter)

	// Act
	res, err := p.Process([]byte("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3318.47\",\"time\":\"2000-01-01T00:00:00Z\"}"))

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, getter.OneDayCalled)
	require.Equal(t, 1, getter.ExistsCalled)
	require.Equal(t, 1, putter.PutCalled)
	require.Equal(t, "{\"id\":\"15887\",\"customer_id\":\"528\",\"accepted\":false}", string(res))
}

func Test_Process_DepositMoreThanLimitPer1Week_Denied(t *testing.T) {
	// Arrange
	getter := mockDepositGetter{
		Res1Week: []*database.Deposit{
			{
				ID:     "1",
				Amount: 1500000,
				Time:   time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			{
				ID:     "2",
				Amount: 400000,
				Time:   time.Date(2000, 0, 3, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	putter := mockDepositPutter{}
	p := New(&getter, &putter)

	// Act
	res, err := p.Process([]byte("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3318.47\",\"time\":\"2000-01-06T00:00:00Z\"}"))

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, getter.OneDayCalled)
	require.Equal(t, "{\"id\":\"15887\",\"customer_id\":\"528\",\"accepted\":false}", string(res))
}

func Test_Process_Success(t *testing.T) {
	// Arrange
	getter := mockDepositGetter{}
	putter := mockDepositPutter{}
	p := New(&getter, &putter)

	// Act
	res, err := p.Process([]byte("{\"id\":\"15887\",\"customer_id\":\"528\",\"load_amount\":\"$3318.47\",\"time\":\"2000-01-06T00:00:00Z\"}"))

	// Assert
	require.NoError(t, err)
	require.Equal(t, 1, getter.OneDayCalled)
	require.Equal(t, 1, putter.PutCalled)
	require.Equal(t, "{\"id\":\"15887\",\"customer_id\":\"528\",\"accepted\":true}", string(res))
}

type mockDepositGetter struct {
	OneDayCalled  int
	OneWeekCalled int
	ExistsCalled  int
	Res1Day       []*database.Deposit
	Res1Week      []*database.Deposit
	ResExists     bool
}

func (m *mockDepositGetter) OneDay(customer_id string, before time.Time) []*database.Deposit {
	m.OneDayCalled++
	return m.Res1Day
}

func (m *mockDepositGetter) OneWeek(customer_id string, before time.Time) []*database.Deposit {
	m.OneWeekCalled++
	return m.Res1Week
}

func (m *mockDepositGetter) Exists(customerID, depositID string) bool {
	m.ExistsCalled++
	return m.ResExists
}

type mockDepositPutter struct {
	PutCalled int
}

func (m *mockDepositPutter) Put(customer_id string, deposit database.Deposit) {
	m.PutCalled++
}
