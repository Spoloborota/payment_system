package model

import "time"

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
}

type Transaction struct {
	ID             uint
	Created        time.Time
	DebitWalletID  uint
	CreditWalletID uint
	Amount         uint
}
