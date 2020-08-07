package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"payment-system/internal/common"
	"payment-system/internal/db"
)

const (
	connPattern   = "%s:%s@tcp(%s:%s)/%s?parseTime=true"
	mysqlType     = "mysql"
	storageType   = "persistence"
	walletIDField = "walletID"
)

type DB struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewDB(host, port, user, pass, db string, logger *zap.Logger) (*DB, error) {
	mysqlLogger := logger.With(zap.String(common.Scope, storageType))
	connStr := fmt.Sprintf(connPattern, user, pass, host, port, db)
	dbx, err := sqlx.Connect(mysqlType, connStr)
	if err != nil {
		mysqlLogger.With(zap.Error(err)).Error("Failed to connect to Mysql db")
		return nil, err
	}
	return &DB{dbx, mysqlLogger}, nil
}

func (d *DB) AddWallet(task AddWalletTask) (uint, error) {
	result, err := d.db.Exec(insertWallet, task.Name, task.Description)
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(addWalletTask, task.String())).Error("Failed to insert wallet to Mysql db")
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(addWalletTask, task.String())).Error("Failed to get new wallet ID after insert it to Mysql db")
		return 0, err
	}

	return uint(id), err
}

func (d *DB) GetWallet(id uint) (db.Wallet, error) {
	var wallet db.Wallet
	err := d.db.Get(&wallet, getWallet, id)
	if err != nil {
		d.logger.With(zap.Error(err), zap.Uint(getWalletTask, id)).Error("Failed to get wallet from Mysql db")
		return db.Wallet{}, err
	}

	return wallet, err
}

func (d *DB) AddTransaction(task AddTransactionTask) (uint64, error) {
	result, err := d.db.Exec(insertTransaction, task.DebitWalletID, task.CreditWalletID, task.Amount, task.Description)
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(addTransactionTask, task.String())).Error("Failed to insert transaction to Mysql db")
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(addTransactionTask, task.String())).Error("Failed to get new transaction ID after insert it to Mysql db")
		return 0, err
	}

	//TODO int64 too small from LastInsertId
	return uint64(id), err
}

func (d *DB) GetTransactionStatus(id uint64) (db.TransactionStatus, string, error) {
	var transaction db.Transaction
	err := d.db.Get(&transaction, getTransactionStatus, id)
	if err != nil {
		d.logger.With(zap.Error(err), zap.Uint64(getTransactionStatusTask, id)).Error("Failed to get transaction status from Mysql db")
		return 0, "", err
	}

	return transaction.Status, transaction.StatusDescription.String, err
}

func (d *DB) GetAllRegisteredTransactions() ([]db.Transaction, error) {
	registeredTransactions := make([]db.Transaction, 0)
	err := d.db.Select(&registeredTransactions, selectAllRegistered)
	if err != nil {
		d.logger.With(zap.Error(err)).Error("Failed to get all registered transactions from Mysql db")
		return nil, err
	}

	return registeredTransactions, err
}

func (d *DB) UpdateTransactionStatus(task UpdateTransactionStatusTask) (err error) {
	_, err = d.db.Exec(updateTransactionStatus, task.Status, task.TransactionID)
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(updateTransactionStatusTask, task.String())).Error("Failed to update transaction Status in Mysql db")
		return
	}
	return
}

func (d *DB) TopUpWalletBalance(task TopUpWalletBalanceTask) error {
	_, err := d.db.Exec(topUpWallet, task.Amount, task.WalletID)
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(topUpWalletBalanceTask, task.String())).Error("Failed to top up wallet in Mysql db")
		return err
	}
	//TODO id only int64 size
	return nil
}

func (d *DB) GetWalletBalance(walletID uint) (uint, error) {
	var balance uint
	err := d.db.Get(&balance, getWalletBalance, walletID)
	if err != nil {
		d.logger.With(zap.Error(err), zap.Uint(walletIDField, walletID)).Error("Failed to get wallet balance from Mysql db")
		return 0, err
	}

	return balance, err
}

func (d *DB) UpdateWalletsBalances(task UpdateWalletsBalancesTask) error {
	tx, err := d.db.Beginx()
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(updateWalletsBalancesTask, task.String())).Error("Failed update wallets balances in Mysql db")
		return err
	}
	_, err = tx.Exec(updateDebitWalletBalance, task.Amount, task.DebitWalletID)
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(updateWalletsBalancesTask, task.String())).Error("Failed to update debit wallet in Mysql db, rolling back")
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			d.logger.With(zap.Error(rollbackErr), zap.String(updateWalletsBalancesTask, task.String())).Error("Failed to rollback update debit wallet in Mysql db")
		}
		return err
	}
	_, err = tx.Exec(updateCreditWalletBalance, task.Amount, task.CreditWalletID)
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(updateWalletsBalancesTask, task.String())).Error("Failed to update credit wallet in Mysql db, rolling back")
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			d.logger.With(zap.Error(rollbackErr), zap.String(updateWalletsBalancesTask, task.String())).Error("Failed to rollback update credit wallet in Mysql db")
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		d.logger.With(zap.Error(err), zap.String(updateWalletsBalancesTask, task.String())).Error("Failed to commit update wallets balances in Mysql db")
		return err
	}
	return nil
}
