package clickhouse

import "fmt"

type Transaction struct {
	ID                                    uint64
	Created                               int64
	DebitWalletID, CreditWalletID, Amount uint32
	Description                           string
	Status                                uint8
}

func (t *Transaction) String() string {
	return fmt.Sprintf("ID: %d, Created: %d, DebitWalletID: %d, CreditWalletID: %d, Amount: %d, Description: %s, Status: %d",
		t.ID, t.Created, t.DebitWalletID, t.CreditWalletID, t.Amount, t.Description, t.Status)
}
