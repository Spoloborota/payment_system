package mysql

import (
	"fmt"
	"payment-system/model"
)

const (
	addTransactionTask          = "AddTransactionTask"
	addWalletTask               = "AddTransactionTask"
	updateTransactionStatusTask = "UpdateTransactionStatusTask"
	topUpWalletBalanceTask      = "TopUpWalletBalanceTask"
	updateWalletsBalancesTask   = "UpdateWalletsBalancesTask"
)

type AddTransactionTask struct {
	DebitWalletID, CreditWalletID, Amount uint
	Description                           string
}

func (a *AddTransactionTask) String() string {
	return fmt.Sprintf("DebitWalletID: %d, CreditWalletID: %d, Amount: %d, Description: %s",
		a.DebitWalletID, a.CreditWalletID, a.Amount, a.Description)
}

type AddWalletTask struct {
	Name, Description string
}

func (a *AddWalletTask) String() string {
	return fmt.Sprintf("Name: %s, Description: %s", a.Name, a.Description)
}

type UpdateTransactionStatusTask struct {
	TransactionID uint64
	Status        model.TransactionStatus
}

func (u *UpdateTransactionStatusTask) String() string {
	return fmt.Sprintf("TransactionID: %d, Status: %d", u.TransactionID, u.Status)
}

type TopUpWalletBalanceTask struct {
	WalletID uint
	Amount   uint
}

func (t *TopUpWalletBalanceTask) String() string {
	return fmt.Sprintf("WalletID: %d, Amount: %d", t.WalletID, t.Amount)
}

type UpdateWalletsBalancesTask struct {
	DebitWalletID, CreditWalletID, Amount uint
}

func (u *UpdateWalletsBalancesTask) String() string {
	return fmt.Sprintf("DebitWalletID: %d, CreditWalletID: %d, Amount: %d", u.DebitWalletID, u.CreditWalletID, u.Amount)
}
