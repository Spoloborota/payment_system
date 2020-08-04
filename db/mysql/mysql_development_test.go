package mysql

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func GetDB(t *testing.T) *DB {
	development, err := zap.NewDevelopment()
	require.NoError(t, err)
	db, err := NewDB("localhost", "3306", "root", "root", "payment_system", development)
	require.NoError(t, err)
	return db
}

func TestAddWallet(t *testing.T) {
	id, err := GetDB(t).AddWallet(AddWalletTask{"XXX", "XXX_DESCR"})
	require.NoError(t, err)
	fmt.Println(id)
}

func TestAddTransaction(t *testing.T) {
	id, err := GetDB(t).AddTransaction(AddTransactionTask{
		DebitWalletID:  3,
		CreditWalletID: 2,
		Amount:         50,
		Description:    "Send 50",
	})
	require.NoError(t, err)
	fmt.Println(id)
}

func TestGetAllRegisteredTransactions(t *testing.T) {
	result, err := GetDB(t).GetAllRegisteredTransactions()
	require.NoError(t, err)
	fmt.Printf("%+v", result)
}

func TestUpdateTransactionStatusTask(t *testing.T) {
	err := GetDB(t).UpdateTransactionStatus(UpdateTransactionStatusTask{1, 2})
	require.NoError(t, err)
}

func TestTopUpWalletBalance(t *testing.T) {
	err := GetDB(t).TopUpWalletBalance(TopUpWalletBalanceTask{2, 10})
	require.NoError(t, err)
}

func TestGetWalletBalance(t *testing.T) {
	val, err := GetDB(t).GetWalletBalance(2)
	require.NoError(t, err)
	fmt.Println(val)
}

func TestUpdateWalletsBalances(t *testing.T) {
	db := GetDB(t)
	err := db.UpdateWalletsBalances(UpdateWalletsBalancesTask{
		DebitWalletID:  2,
		CreditWalletID: 3,
		Amount:         10,
	})
	require.NoError(t, err)
	err = db.UpdateWalletsBalances(UpdateWalletsBalancesTask{
		DebitWalletID:  2,
		CreditWalletID: 3,
		Amount:         100000,
	})
	require.Error(t, err)
}
