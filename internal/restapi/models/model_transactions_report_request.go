package models

type TransactionReportRequest struct {
	IsTopUps       bool
	CreatedFrom    int64
	CreatedTo      int64
	DebitWalletID  uint32
	CreditWalletID uint32
}
