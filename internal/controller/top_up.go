package controller

import (
	"fmt"
	"go.uber.org/zap"
	"payment-system/internal/db"
	"payment-system/internal/db/clickhouse"
	"payment-system/internal/db/mysql"
	"payment-system/internal/restapi/models"
	"time"
)

const (
	bankWalletID          = 1
	transactionIDForTopUp = 0
)

func (c *Controller) TopupWallet(twr models.TopupWalletRequest) error {
	err := c.mysqlDB.TopUpWalletBalance(mysql.TopUpWalletBalanceTask{
		WalletID: twr.Id,
		Amount:   twr.Amount,
	})
	if err != nil {
		return fmt.Errorf("failed to topup wallet: %w", err)
	}

	go func(twr models.TopupWalletRequest) {
		err := c.clickHouseDB.AddTransaction(clickhouse.Transaction{
			ID:                transactionIDForTopUp,
			Created:           time.Now().Unix(),
			DebitWalletID:     uint32(twr.Id),
			CreditWalletID:    bankWalletID,
			Amount:            int32(twr.Amount),
			Description:       "Top Up Balance operation",
			Status:            uint8(db.Success),
			StatusDescription: "Success",
		})
		if err != nil {
			c.logger.Error("Failed to put transaction to journal db", zap.Error(err))
		}
	}(twr)

	return nil
}
