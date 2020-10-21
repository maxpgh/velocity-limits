package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	// Act
	db := New()

	// Assert
	require.NotNil(t, db)
}
