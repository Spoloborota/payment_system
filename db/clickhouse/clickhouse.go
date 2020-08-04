package clickhouse

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"payment-system/internal"
)

const (
	connPattern    = "tcp://%s:%s?database=%s"
	clickhouseType = "clickhouse"
	storageType    = "journal"

	addTransaction = "addTransaction"
)

type DB struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewDB(host, port, db string, logger *zap.Logger) (*DB, error) {
	clickhouseLogger := logger.With(zap.String(internal.Scope, storageType))
	connStr := fmt.Sprintf(connPattern, host, port, db)
	dbx, err := sqlx.Connect(clickhouseType, connStr)
	if err != nil {
		clickhouseLogger.With(zap.Error(err)).Error("Failed to connect to ClickHouse db")
		return nil, err
	}
	return &DB{dbx, clickhouseLogger}, nil
}

func (db *DB) AddTransaction(tx Transaction) error {
	ttx, err := db.db.Beginx()
	if err != nil {
		db.logger.With(zap.Error(err), zap.String(addTransaction, tx.String())).Error("Failed to insert transaction to Clickhouse db")
		return err
	}
	_, err = ttx.Exec(insertTransaction, tx.ID, tx.Created, tx.DebitWalletID, tx.CreditWalletID, tx.Amount, tx.Description, tx.Status)
	if err != nil {
		db.logger.With(zap.Error(err), zap.String(addTransaction, tx.String())).Error("Failed to insert transaction to Clickhouse db")
		return err
	}
	err = ttx.Commit()
	if err != nil {
		db.logger.With(zap.Error(err), zap.String(addTransaction, tx.String())).Error("Failed to insert transaction to Clickhouse db")
		return err
	}

	return nil
}
