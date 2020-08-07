package clickhouse

import (
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"payment-system/internal/common"
	"strconv"
	"strings"
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
	clickhouseLogger := logger.With(zap.String(common.Scope, storageType))
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
	_, err = ttx.Exec(insertTransaction, tx.ID, tx.Created, tx.DebitWalletID, tx.CreditWalletID, tx.Amount, tx.Description, tx.Status, tx.StatusDescription)
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

//TODO add pagination
func (db *DB) TransactionsReport(query ReportQuery) ([]Transaction, error) {
	var sb strings.Builder
	var queried bool
	var isWhereAdded bool
	sb.WriteString(getReport)
	if query.CreatedFrom != 0 {
		if queried {
			sb.WriteString(" AND")
		}
		queried = true
		if !isWhereAdded {
			sb.WriteString(" WHERE")
		}
		isWhereAdded = true
		sb.WriteString(" created >= ")
		sb.WriteString(strconv.FormatInt(query.CreatedFrom, 10))
	}
	if query.CreatedTo != 0 {
		if queried {
			sb.WriteString(" AND")
		}
		queried = true
		if !isWhereAdded {
			sb.WriteString(" WHERE")
		}
		isWhereAdded = true
		sb.WriteString(" created <= ")
		sb.WriteString(strconv.FormatInt(query.CreatedTo, 10))
	}
	if query.DebitWalletID != 0 {
		if queried {
			sb.WriteString(" AND")
		}
		queried = true
		if !isWhereAdded {
			sb.WriteString(" WHERE")
		}
		isWhereAdded = true
		sb.WriteString(" debit_wallet_id = ")
		sb.WriteString(strconv.FormatUint(uint64(query.DebitWalletID), 10))
	}
	if query.IsTopUps {
		if queried {
			sb.WriteString(" AND")
		}
		queried = true
		if !isWhereAdded {
			sb.WriteString(" WHERE")
		}
		isWhereAdded = true
		sb.WriteString(" credit_wallet_id = 1")
	} else {
		if query.CreditWalletID != 0 {
			if queried {
				sb.WriteString(" AND")
			}
			queried = true
			if !isWhereAdded {
				sb.WriteString(" WHERE")
			}
			isWhereAdded = true
			sb.WriteString(" credit_wallet_id = ")
			sb.WriteString(strconv.FormatUint(uint64(query.CreditWalletID), 10))
		}
	}
	sb.WriteString(";")

	resultQuery := sb.String()
	ttx, err := db.db.Beginx()
	if err != nil {
		db.logger.With(zap.Error(err), zap.String(addTransaction, query.String())).Error("Failed to start query to Clickhouse db")
		return nil, err
	}
	transactions := make([]Transaction, 0)
	err = ttx.Select(&transactions, resultQuery)
	if err != nil {
		db.logger.With(zap.Error(err), zap.String(addTransaction, query.String())).Error("Failed to start query to Clickhouse db")
		return nil, err
	}
	err = ttx.Commit()
	if err != nil {
		db.logger.With(zap.Error(err), zap.String(addTransaction, query.String())).Error("Failed to start query to Clickhouse db")
		return nil, err
	}

	return transactions, nil
}
