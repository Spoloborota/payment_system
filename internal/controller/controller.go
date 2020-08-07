package controller

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"payment-system/internal/db"
	"payment-system/internal/db/clickhouse"
	"payment-system/internal/db/mysql"
	"time"
)

type Controller struct {
	logger       *zap.Logger
	mysqlDB      *mysql.DB
	clickHouseDB *clickhouse.DB
}

func NewController(logger *zap.Logger, mysqlDB *mysql.DB, clickHouseDB *clickhouse.DB) *Controller {
	return &Controller{
		logger:       logger,
		mysqlDB:      mysqlDB,
		clickHouseDB: clickHouseDB,
	}
}

func (c *Controller) StartTimerTask() chan<- bool {
	interval := viper.GetInt("processing_task_timeout_sec")
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-ticker.C:
				c.ProcessRegisteredTransactions()
			}
		}
	}()
	return quit
}

func (c *Controller) ProcessRegisteredTransactions() {
	registeredTransactions, err := c.mysqlDB.GetAllRegisteredTransactions()
	if err != nil {
		c.logger.Error("Failed to receive all registered tasks", zap.Error(err))
		return
	}
	if len(registeredTransactions) == 0 {
		return
	}
	for _, tx := range registeredTransactions {
		wallet, err := c.mysqlDB.GetWallet(tx.CreditWalletID)
		if err != nil {
			c.logger.Error("Failed to get credit walet for registered task", zap.Error(err), zap.Uint("credit_wallet", tx.CreditWalletID))
			c.UpdateTransactions(tx, db.Failed, "Failed to get credit wallet from db")
			continue
		}
		_, err = c.mysqlDB.GetWallet(tx.DebitWalletID)
		if err != nil {
			c.logger.Error("Failed to get debit wallet for registered task", zap.Error(err), zap.Uint("debit_wallet", tx.DebitWalletID))
			c.UpdateTransactions(tx, db.Failed, "Failed to get debit wallet from db")
			continue
		}
		if wallet.Balance >= tx.Amount {
			err := c.mysqlDB.UpdateWalletsBalances(mysql.UpdateWalletsBalancesTask{
				DebitWalletID:  tx.DebitWalletID,
				CreditWalletID: tx.CreditWalletID,
				Amount:         tx.Amount,
			})
			if err != nil {
				c.logger.Error("Failed to update balances for registered task", zap.Error(err), zap.Uint("credit_wallet", tx.CreditWalletID))
				c.UpdateTransactions(tx, db.Failed, "Failed to update wallet balance")
				continue
			}
			c.UpdateTransactions(tx, db.Success, "Success")
		} else {
			c.logger.Error("Not enough money for credit from wallet for task", zap.Error(err), zap.Uint("credit_wallet", tx.CreditWalletID), zap.Uint64("transaction_id", tx.ID))
			c.UpdateTransactions(tx, db.Failed, fmt.Sprintf("Not enough money on credit wallet: %d", wallet.Balance))
		}
	}
}

//TODO retry on fail and check rows
func (c *Controller) UpdateTransactions(tx db.Transaction, status db.TransactionStatus, description string) {
	err := c.mysqlDB.UpdateTransactionStatus(mysql.UpdateTransactionStatusTask{
		TransactionID: tx.ID,
		Status:        status,
	})
	if err != nil {
		c.logger.Error("Failed to update status of transaction in mysql", zap.Error(err), zap.Uint64("transaction_id", tx.ID))
	}

	err = c.clickHouseDB.AddTransaction(clickhouse.Transaction{
		ID:                tx.ID,
		Created:           tx.Created.Unix(),
		Amount:            int32(tx.Amount),
		DebitWalletID:     uint32(tx.DebitWalletID),
		CreditWalletID:    uint32(tx.CreditWalletID),
		Description:       tx.Description,
		Status:            uint8(status),
		StatusDescription: description,
	})
	if err != nil {
		c.logger.Error("Failed to insert positive transaction into clickhouse", zap.Error(err), zap.Uint64("transaction_id", tx.ID))
	}

	err = c.clickHouseDB.AddTransaction(clickhouse.Transaction{
		ID:                tx.ID,
		Created:           tx.Created.Unix(),
		Amount:            int32(-tx.Amount),
		DebitWalletID:     uint32(tx.CreditWalletID),
		CreditWalletID:    uint32(tx.DebitWalletID),
		Description:       tx.Description,
		Status:            uint8(status),
		StatusDescription: description,
	})
	if err != nil {
		c.logger.Error("Failed to insert negative transaction into clickhouse", zap.Error(err), zap.Uint64("transaction_id", tx.ID))
	}
}
