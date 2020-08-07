package controller

import (
	"fmt"
	"payment-system/internal/db/mysql"
	"payment-system/internal/restapi/models"
)

func (c *Controller) AddTransaction(ctr models.CreateTransactionRequest) (uint64, error) {
	id, err := c.mysqlDB.AddTransaction(mysql.AddTransactionTask{
		DebitWalletID:  ctr.DebitId,
		CreditWalletID: ctr.CreditId,
		Amount:         ctr.Amount,
		Description:    ctr.Description,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create transaction: %w", err)
	}

	return id, nil
}
