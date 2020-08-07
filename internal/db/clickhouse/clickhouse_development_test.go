//These tests are used only in development or manual testing. They are not configured for unit or integration testing.

package clickhouse

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func GetDB(t *testing.T) *DB {
	development, err := zap.NewDevelopment()
	require.NoError(t, err)
	db, err := NewDB("localhost", "9000", "payment_system", development)
	require.NoError(t, err)
	return db
}

func TestAddTransaction(t *testing.T) {
	db := GetDB(t)
	err := db.AddTransaction(Transaction{
		ID:             2,
		Created:        time.Now().Unix(),
		DebitWalletID:  2,
		CreditWalletID: 3,
		Amount:         200,
		Description:    "XXX",
		Status:         2,
	})
	require.NoError(t, err)
}
