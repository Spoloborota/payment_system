package clickhouse

import (
	"fmt"
)

type Transaction struct {
	ID                uint64
	Created           int64
	Amount            int32
	DebitWalletID     uint32 `db:"debit_wallet_id"`
	CreditWalletID    uint32 `db:"credit_wallet_id"`
	Description       string
	Status            uint8
	StatusDescription string `db:"status_description"`
}

func (t *Transaction) String() string {
	return fmt.Sprintf("ID: %d, Created: %d, DebitWalletID: %d, CreditWalletID: %d, Amount: %d, Description: %s, Status: %d, StatusDescription: %s",
		t.ID, t.Created, t.DebitWalletID, t.CreditWalletID, t.Amount, t.Description, t.Status, t.StatusDescription)
}

type ReportQuery struct {
	IsTopUps       bool
	CreatedFrom    int64
	CreatedTo      int64
	DebitWalletID  uint32
	CreditWalletID uint32
}

func (r *ReportQuery) String() string {
	return fmt.Sprintf("IsTopUps: %t, CreatedFrom: %d, CreatedTo: %d, DebitWalletID: %d, CreditWalletID: %d",
		r.IsTopUps, r.CreatedFrom, r.CreatedTo, r.CreditWalletID, r.CreditWalletID)
}
