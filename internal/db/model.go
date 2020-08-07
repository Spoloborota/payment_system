package db

import (
	"database/sql"
	"time"
)

const (
	Registered TransactionStatus = iota
	Success
	Failed
)

type TransactionStatus uint8

type WalletType struct {
	ID          uint
	Name        string
	Description string
}

type Wallet struct {
	ID           uint
	Name         string
	Description  string
	WalletTypeID uint
	Balance      uint
}

type Transaction struct {
	ID                uint64
	Created           time.Time
	DebitWalletID     uint `db:"debit_wallet_id"`
	CreditWalletID    uint `db:"credit_wallet_id"`
	Amount            uint
	Description       string
	Status            TransactionStatus
	StatusDescription sql.NullString `db:"status_description"`
}
